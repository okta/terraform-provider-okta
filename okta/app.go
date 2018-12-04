package okta

import (
	"fmt"
	"sync"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

var baseAppSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "name of app.",
	},
	"label": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Pretty name of app.",
	},
	"sign_on_mode": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Sign on mode of application.",
	},
	"users": &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"scope": &schema.Schema{
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Scope of application user.",
					},
					"id": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "User ID.",
					},
					"username": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Username for user.",
					},
					"password": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Password for user application.",
					},
				},
			},
		},
		Description: "List of users associated with the application",
	},
	"groups": &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "List of groups associated with the application",
	},
	"status": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "ACTIVE",
		ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
		Description:  "Status of application.",
	},
	"app_settings": {
		Type:        schema.TypeMap,
		Optional:    true,
		Description: "Application settings",
		Elem:        schema.TypeString,
	},
}

var baseSwaAppSchema = map[string]*schema.Schema{
	"accessibility_self_service": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Enable self service",
	},
	"accessibility_error_redirect_url": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "Custom error page URL",
		ValidateFunc: validateIsURL,
	},
	"auto_submit_toolbar": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Display auto submit toolbar",
	},
	"hide_ios": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Do not display application icon on mobile app",
	},
	"hide_web": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Do not display application icon to users",
	},
	"user_name_template": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Username template",
	},
	"user_name_template_type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Username template type",
	},
}

// Wish there was some better polymorphism that could make these similarities easier to deal with
func appRead(d *schema.ResourceData, name, status, signOn, label string, accy *okta.ApplicationAccessibility, vis *okta.ApplicationVisibility) {
	d.Set("name", name)
	d.Set("status", status)
	d.Set("sign_on_mode", signOn)
	d.Set("label", label)
	d.Set("accessibility_self_service", accy.SelfService)
	d.Set("accessibility_error_redirect_url", accy.ErrorRedirectUrl)
	d.Set("auto_submit_toolbar", vis.AutoSubmitToolbar)
	d.Set("hide_ios", vis.Hide.IOS)
	d.Set("hide_web", vis.Hide.Web)
}

func buildAppSchema(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseAppSchema, appSchema)
}

func buildSchemeCreds(d *schema.ResourceData) *okta.SchemeApplicationCredentials {
	revealPass := d.Get("reveal_password").(bool)

	return &okta.SchemeApplicationCredentials{
		RevealPassword: &revealPass,
		Scheme:         d.Get("credentials_scheme").(string),
		UserName:       d.Get("shared_username").(string),
		Password: &okta.PasswordCredential{
			Value: d.Get("shared_password").(string),
		},
	}
}

func buildSwaAppSchema(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	s := buildAppSchema(appSchema)
	return buildSchema(baseSwaAppSchema, s)
}

func buildVisibility(d *schema.ResourceData) *okta.ApplicationVisibility {
	autoSubmit := d.Get("auto_submit_toolbar").(bool)
	hideMobile := d.Get("hide_ios").(bool)
	hideWeb := d.Get("hide_web").(bool)
	return &okta.ApplicationVisibility{
		AutoSubmitToolbar: &autoSubmit,
		Hide: &okta.ApplicationVisibilityHide{
			IOS: &hideMobile,
			Web: &hideWeb,
		},
	}
}

func fetchApp(d *schema.ResourceData, m interface{}, app okta.App) error {
	client := getOktaClientFromMetadata(m)
	params := &query.Params{}
	_, response, err := client.Application.GetApplication(d.Id(), app, params)
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		app = nil
		return nil
	}

	return err
}

func handleAppGroups(id string, d *schema.ResourceData, client *okta.Client) []func() error {
	existingGroup, _, _ := client.Application.ListApplicationGroupAssignments(id, &query.Params{})
	var asyncActionList []func() error
	rawArr, ok := d.Get("groups").([]interface{})

	if ok {
		for _, thing := range rawArr {
			g := thing.(string)
			contains := false

			for _, eGroup := range existingGroup {
				if eGroup.Id == g {
					contains = true
					break
				}
			}

			if !contains {
				asyncActionList = append(asyncActionList, func() error {
					_, _, err := client.Application.CreateApplicationGroupAssignment(id, g, okta.ApplicationGroupAssignment{})
					return err
				})
			}
		}
	}

	for _, group := range existingGroup {
		contains := false
		for _, thing := range rawArr {
			g := thing.(string)
			if g == group.Id {
				contains = true
				break
			}
		}

		if !contains {
			asyncActionList = append(asyncActionList, func() error {
				_, err := client.Application.DeleteApplicationGroupAssignment(id, group.Id)
				return err
			})
		}
	}

	return asyncActionList
}

// Handles the assigning of groups and users to Applications. Does so asynchronously.
func handleAppGroupsAndUsers(id string, d *schema.ResourceData, m interface{}) error {
	var wg sync.WaitGroup
	resultChan := make(chan []*result, 1)
	client := getOktaClientFromMetadata(m)

	groupHandlers := handleAppGroups(id, d, client)
	userHandlers := handleAppUsers(id, d, client)
	con := getParallelismFromMetadata(m)
	promiseAll(con, &wg, resultChan, append(groupHandlers, userHandlers...)...)
	wg.Wait()

	return getPromiseError(<-resultChan, "failed to associate user or groups with application")
}

func handleAppUsers(id string, d *schema.ResourceData, client *okta.Client) []func() error {
	// Looking upstream for existing user's, rather then the config for accuracy.
	existingUsers, _, _ := client.Application.ListApplicationUsers(id, &query.Params{})
	var asyncActionList []func() error
	rawList, ok := d.GetOk("users")
	var userIDList []string

	if ok {
		userLen := len(rawList.([]interface{}))
		userIDList = make([]string, userLen)

		for i := 0; i < userLen; i++ {
			path := fmt.Sprintf("users.%v", i)
			uID := d.Get(fmt.Sprintf("%s.id", path)).(string)
			userIDList[i] = uID
			contains := false

			for _, u := range existingUsers {
				if u.Id == uID && u.Scope == "USER" {
					contains = true
					break
				}
			}

			if !contains {
				asyncActionList = append(asyncActionList, func() error {
					_, _, err := client.Application.AssignUserToApplication(id, okta.AppUser{
						Id: uID,
						Credentials: &okta.AppUserCredentials{
							UserName: d.Get(fmt.Sprintf("%s.username", path)).(string),
							Password: &okta.AppUserPasswordCredential{
								Value: d.Get(fmt.Sprintf("%s.password", path)).(string),
							},
						},
					})

					return err
				})
			}
		}

	}

	for _, user := range existingUsers {
		if user.Scope == "USER" {
			contains := false
			for _, uID := range userIDList {
				if uID == user.Id {
					contains = true
					break
				}
			}

			if !contains {
				asyncActionList = append(asyncActionList, func() error {
					_, err := client.Application.DeleteApplicationUser(id, user.Id)

					return err
				})
			}
		}
	}

	return asyncActionList
}

func resourceAppExists(d *schema.ResourceData, m interface{}) (bool, error) {
	app := okta.NewApplication()
	err := fetchApp(d, m, app)

	// Not sure if a non-nil app with an empty ID is possible but checking to avoid false positives.
	return app != nil && app.Id != "", err
}

func setAppStatus(d *schema.ResourceData, client *okta.Client, status string, desiredStatus string) error {
	var err error
	if status != desiredStatus {
		if desiredStatus == "INACTIVE" {
			_, err = client.Application.DeactivateApplication(d.Id())
		} else if desiredStatus == "ACTIVE" {
			_, err = client.Application.ActivateApplication(d.Id())
		}
	}

	return err
}

func syncGroupsAndUsers(id string, d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	userList, _, err := client.Application.ListApplicationUsers(id, &query.Params{})
	if err != nil {
		return err
	}

	groupList, _, err := client.Application.ListApplicationGroupAssignments(id, &query.Params{})
	if err != nil {
		return err
	}
	flatGroupList := make([]string, len(groupList))

	for i, g := range groupList {
		flatGroupList[i] = g.Id
	}

	d.Set("groups", flatGroupList)
	var flattenedUserList []interface{}

	for _, user := range userList {
		if user.Scope == "USER" {
			flattenedUserList = append(flattenedUserList, map[string]interface{}{
				"id":       user.Id,
				"username": user.Credentials.UserName,
			})
		}
	}
	flatMap := map[string]interface{}{}

	if len(flattenedUserList) > 0 {
		flatMap["users"] = flattenedUserList
	}

	if len(flatGroupList) > 0 {
		flatMap["groups"] = flatGroupList
	}

	return setNonPrimitives(d, flatMap)
}
