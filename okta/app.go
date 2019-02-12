package okta

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
	"github.com/peterhellberg/link"
)

type appID struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type appFilters struct {
	ID                string
	Label             string
	LabelPrefix       string
	ShortCircuitCount int
}

func (f *appFilters) shouldShortCircuit(appList []*appID) bool {
	if f.LabelPrefix != "" {
		return false
	}

	if f.ID != "" && f.Label != "" {
		return len(appList) > 1
	}

	if f.ID != "" || f.Label != "" {
		return len(appList) > 0
	}

	return false
}

var appUserResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"scope": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Scope of application user.",
		},
		"id": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User ID.",
		},
		"username": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Username for user.",
		},
		"password": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Password for user application.",
		},
	},
}

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
		Type:        schema.TypeSet,
		Optional:    true,
		Elem:        appUserResource,
		Description: "Users associated with the application",
	},
	"groups": &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Groups associated with the application",
	},
	"status": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "ACTIVE",
		ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
		Description:  "Status of application.",
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
	var (
		asyncActionList []func() error
		rawArr          []interface{}
	)

	if arr, ok := d.GetOk("groups"); ok {
		rawArr = arr.(*schema.Set).List()
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
				return suppressErrorOn404(client.Application.DeleteApplicationGroupAssignment(id, group.Id))
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
	var (
		asyncActionList []func() error
		users           []interface{}
		userIDList      []string
	)

	if set, ok := d.GetOk("users"); ok {
		users = set.(*schema.Set).List()
		userIDList = make([]string, len(users))

		for i, user := range users {
			u := user.(map[string]interface{})
			uID := u["id"].(string)
			userIDList[i] = uID
			contains := false

			for _, u := range existingUsers {
				if u.Id == uID && u.Scope == "USER" {
					contains = true
					break
				}
			}

			if !contains {
				username := u["username"].(string)
				// Not required
				password, _ := u["password"].(string)

				asyncActionList = append(asyncActionList, func() error {
					_, _, err := client.Application.AssignUserToApplication(id, okta.AppUser{
						Id: uID,
						Credentials: &okta.AppUserCredentials{
							UserName: username,
							Password: &okta.AppUserPasswordCredential{
								Value: password,
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
					return suppressErrorOn404(client.Application.DeleteApplicationUser(id, user.Id))
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
	flatGroupList := make([]interface{}, len(groupList))

	for i, g := range groupList {
		flatGroupList[i] = g.Id
	}

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
		flatMap["users"] = schema.NewSet(schema.HashResource(appUserResource), flattenedUserList)
	}

	if len(flatGroupList) > 0 {
		flatMap["groups"] = schema.NewSet(schema.HashString, flatGroupList)
	}

	return setNonPrimitives(d, flatMap)
}

// setAppSettings available preconfigured SAML and OAuth applications vary wildly on potential app settings, thus
// it is a generic map. This logic simply weeds out any empty string values.
func setAppSettings(d *schema.ResourceData, settings *okta.ApplicationSettingsApplication) error {
	flatMap := map[string]interface{}{}

	for key, val := range *settings {
		if str, ok := val.(string); ok {
			if str != "" {
				flatMap[key] = str
			}
		} else if val != nil {
			flatMap[key] = val
		}
	}
	payload, err := json.Marshal(flatMap)
	if err != nil {
		return err
	}

	return d.Set("app_settings_json", string(payload))
}

func listApps(c *Config, filters *appFilters) ([]*appID, error) {
	return collectApps([]*appID{}, &query.Params{}, c.supplementClient.requestExecutor, filters)
}

// Recursively list apps until no next links are returned
func collectApps(apps []*appID, qp *query.Params, requestExecutor *okta.RequestExecutor, filters *appFilters) ([]*appID, error) {
	req, err := requestExecutor.NewRequest("GET", "/api/v1/apps", nil)
	if err != nil {
		return nil, err
	}

	var appList []*appID
	res, err := requestExecutor.Do(req, &appList)
	if err != nil {
		return nil, err
	}

	apps = append(apps, filterApp(appList, filters)...)
	linkList := link.Parse(res.Header.Get("link"))
	if linkList == nil {
		return apps, nil
	}
	if filters.shouldShortCircuit(apps) {
		return apps, nil
	}

	for _, l := range linkList {
		if l.Rel == "next" {
			parsedURL, err := url.Parse(l.URI)
			if err != nil {
				continue
			}
			q := parsedURL.Query()
			qp.After = q.Get("after")
			limit, err := strconv.ParseInt(q.Get("limit"), 10, 64)
			if err == nil {
				qp.Limit = limit
			}
			return collectApps(apps, qp, requestExecutor, filters)
		}
	}
	return apps, err
}

func filterApp(appList []*appID, filter *appFilters) []*appID {
	filteredList := []*appID{}
	for _, app := range appList {
		if filter.ID == app.ID {
			filteredList = append(filteredList, app)
		}

		if filter.Label == app.Label {
			filteredList = append(filteredList, app)
		}

		if filter.LabelPrefix != "" && strings.HasPrefix(app.Label, filter.LabelPrefix) {
			filteredList = append(filteredList, app)
		}

		if filter.Label == "" && filter.ID == "" && filter.LabelPrefix == "" {
			filteredList = append(filteredList, app)
		}
	}
	return filteredList
}
