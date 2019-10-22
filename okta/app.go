package okta

import (
	"encoding/json"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

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

var appVisibilitySchema = map[string]*schema.Schema{
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
}

var baseappSwaSchema = map[string]*schema.Schema{
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

func buildAppSchemaWithVisibility(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	schema := buildSchema(baseAppSchema, appSchema)
	return buildSchema(appVisibilitySchema, schema)
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

func buildAppSwaSchema(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	s := buildAppSchema(appSchema)
	return buildSchema(baseappSwaSchema, s)
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
	return fetchAppById(d.Id(), m, app)
}

func fetchAppById(id string, m interface{}, app okta.App) error {
	client := getOktaClientFromMetadata(m)
	_, response, err := client.Application.GetApplication(id, app, nil)
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		app = nil
		return nil
	}

	return responseErr(response, err)
}

func updateAppById(id string, m interface{}, app okta.App) error {
	client := getOktaClientFromMetadata(m)
	_, response, err := client.Application.UpdateApplication(id, app)
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		app = nil
		return nil
	}

	return responseErr(response, err)
}

func handleAppGroups(id string, d *schema.ResourceData, client *okta.Client) []func() error {
	existingGroup, _, _ := client.Application.ListApplicationGroupAssignments(id, &query.Params{})
	var (
		asyncActionList []func() error
		groupIdList     []string
	)

	if arr, ok := d.GetOk("groups"); ok {
		rawArr := arr.(*schema.Set).List()
		groupIdList = make([]string, len(rawArr))

		for i, gID := range rawArr {
			groupID := gID.(string)
			groupIdList[i] = groupID

			if !containsGroup(existingGroup, groupID) {
				asyncActionList = append(asyncActionList, func() error {
					_, resp, err := client.Application.CreateApplicationGroupAssignment(id, groupID, okta.ApplicationGroupAssignment{})
					return responseErr(resp, err)
				})
			}
		}
	}

	for _, group := range existingGroup {
		if !contains(groupIdList, group.Id) {
			groupID := group.Id
			asyncActionList = append(asyncActionList, func() error {
				return suppressErrorOn404(client.Application.DeleteApplicationGroupAssignment(id, groupID))
			})
		}
	}

	return asyncActionList
}

func containsGroup(groupList []*okta.ApplicationGroupAssignment, id string) bool {
	for _, group := range groupList {
		if group.Id == id {
			return true
		}
	}
	return false
}

func containsAppUser(userList []*okta.AppUser, id string) bool {
	for _, user := range userList {
		if user.Id == id && user.Scope == "USER" {
			return true
		}
	}
	return false
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
			userProfile := user.(map[string]interface{})
			uID := userProfile["id"].(string)
			userIDList[i] = uID

			if !containsAppUser(existingUsers, uID) {
				username := userProfile["username"].(string)
				// Not required
				password, _ := userProfile["password"].(string)

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
			if !contains(userIDList, user.Id) {
				userID := user.Id
				asyncActionList = append(asyncActionList, func() error {
					return suppressErrorOn404(client.Application.DeleteApplicationUser(id, userID, nil))
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
	if status != desiredStatus {
		if desiredStatus == "INACTIVE" {
			return responseErr(client.Application.DeactivateApplication(d.Id()))
		} else if desiredStatus == "ACTIVE" {
			return responseErr(client.Application.ActivateApplication(d.Id()))
		}
	}

	return nil
}

func syncGroupsAndUsers(id string, d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	// Temporary high limit to avoid issues short term. Need to support pagination here
	userList, _, err := client.Application.ListApplicationUsers(id, &query.Params{Limit: 200})
	if err != nil {
		return err
	}

	// Temporary high limit to avoid issues short term. Need to support pagination here
	groupList, _, err := client.Application.ListApplicationGroupAssignments(id, &query.Params{Limit: 200})
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

func syncSamlSettings(d *schema.ResourceData, set *okta.SamlApplicationSettings) error {
	d.Set("default_relay_state", set.SignOn.DefaultRelayState)
	d.Set("sso_url", set.SignOn.SsoAcsUrl)
	d.Set("recipient", set.SignOn.Recipient)
	d.Set("destination", set.SignOn.Destination)
	d.Set("audience", set.SignOn.Audience)
	d.Set("idp_issuer", set.SignOn.IdpIssuer)
	d.Set("subject_name_id_template", set.SignOn.SubjectNameIdTemplate)
	d.Set("subject_name_id_format", set.SignOn.SubjectNameIdFormat)
	d.Set("response_signed", set.SignOn.ResponseSigned)
	d.Set("assertion_signed", set.SignOn.AssertionSigned)
	d.Set("signature_algorithm", set.SignOn.SignatureAlgorithm)
	d.Set("digest_algorithm", set.SignOn.DigestAlgorithm)
	d.Set("honor_force_authn", set.SignOn.HonorForceAuthn)
	d.Set("authn_context_class_ref", set.SignOn.AuthnContextClassRef)

	attrStatements := set.SignOn.AttributeStatements
	arr := make([]map[string]interface{}, len(attrStatements))

	for i, st := range attrStatements {
		arr[i] = map[string]interface{}{
			"name":         st.Name,
			"namespace":    st.Namespace,
			"type":         st.Type,
			"values":       st.Values,
			"filter_type":  st.FilterType,
			"filter_value": st.FilterValue,
		}
	}

	return setNonPrimitives(d, map[string]interface{}{
		"attribute_statements": arr,
	})
}
