package okta

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

var appUserResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"scope": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Scope of application user.",
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User ID.",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Username for user.",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Password for user application.",
		},
	},
}

var baseAppSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "name of app.",
	},
	"label": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Pretty name of app.",
	},
	"sign_on_mode": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Sign on mode of application.",
	},
	"users": {
		Type:        schema.TypeSet,
		Optional:    true,
		Elem:        appUserResource,
		Description: "Users associated with the application",
	},
	"groups": {
		Type:        schema.TypeSet,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Groups associated with the application",
	},
	"status": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          statusActive,
		ValidateDiagFunc: stringInSlice([]string{statusActive, statusInactive}),
		Description:      "Status of application.",
	},
}

var appVisibilitySchema = map[string]*schema.Schema{
	"auto_submit_toolbar": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Display auto submit toolbar",
	},
	"hide_ios": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Do not display application icon on mobile app",
	},
	"hide_web": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Do not display application icon to users",
	},
}

var baseAppSwaSchema = map[string]*schema.Schema{
	"accessibility_self_service": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Enable self service",
	},
	"accessibility_error_redirect_url": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Custom error page URL",
		ValidateDiagFunc: stringIsURL(validURLSchemes...),
	},
	"auto_submit_toolbar": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Display auto submit toolbar",
	},
	"hide_ios": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Do not display application icon on mobile app",
	},
	"hide_web": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Do not display application icon to users",
	},
	"user_name_template": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "${source.login}",
		Description: "Username template",
	},
	"user_name_template_suffix": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Username template suffix",
	},
	"user_name_template_type": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          "BUILT_IN",
		Description:      "Username template type",
		ValidateDiagFunc: stringInSlice([]string{"NONE", "CUSTOM", "BUILT_IN"}),
	},
}

var attributeStatements = map[string]*schema.Schema{
	"filter_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Type of group attribute filter",
	},
	"filter_value": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filter value to use",
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"namespace": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified",
		ValidateDiagFunc: stringInSlice([]string{
			"urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified",
			"urn:oasis:names:tc:SAML:2.0:attrname-format:uri",
			"urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
		}),
	},
	"type": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "EXPRESSION",
		ValidateDiagFunc: stringInSlice([]string{
			"EXPRESSION",
			"GROUP",
		}),
	},
	"values": {
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
}

var appSamlDiffSuppressFunc = func(k, old, new string, d *schema.ResourceData) bool {
	// Conditional default
	return new == "" && old == "http://www.okta.com/${org.externalKey}"
}

// Wish there was some better polymorphism that could make these similarities easier to deal with
func appRead(d *schema.ResourceData, name, status, signOn, label string, accy *okta.ApplicationAccessibility, vis *okta.ApplicationVisibility) {
	_ = d.Set("name", name)
	_ = d.Set("status", status)
	_ = d.Set("sign_on_mode", signOn)
	_ = d.Set("label", label)
	_ = d.Set("accessibility_self_service", *accy.SelfService)
	_ = d.Set("accessibility_error_redirect_url", accy.ErrorRedirectUrl)
	_ = d.Set("auto_submit_toolbar", vis.AutoSubmitToolbar)
	_ = d.Set("hide_ios", vis.Hide.IOS)
	_ = d.Set("hide_web", vis.Hide.Web)
}

func buildAppSchema(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseAppSchema, appSchema)
}

func buildAppSchemaWithVisibility(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseAppSchema, appVisibilitySchema, appSchema)
}

func buildSchemeCreds(d *schema.ResourceData) *okta.SchemeApplicationCredentials {
	revealPass := d.Get("reveal_password").(bool)

	return &okta.SchemeApplicationCredentials{
		RevealPassword: &revealPass,
		Scheme:         d.Get("credentials_scheme").(string),
		UserNameTemplate: &okta.ApplicationCredentialsUsernameTemplate{
			Template: d.Get("user_name_template").(string),
			Type:     d.Get("user_name_template_type").(string),
			Suffix:   d.Get("user_name_template_suffix").(string),
		},
		UserName: d.Get("shared_username").(string),
		Password: &okta.PasswordCredential{
			Value: d.Get("shared_password").(string),
		},
	}
}

func buildAppSwaSchema(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseAppSchema, baseAppSwaSchema, appSchema)
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

func fetchApp(ctx context.Context, d *schema.ResourceData, m interface{}, app okta.App) error {
	return fetchAppByID(ctx, d.Id(), m, app)
}

func fetchAppByID(ctx context.Context, id string, m interface{}, app okta.App) error {
	_, resp, err := getOktaClientFromMetadata(m).Application.GetApplication(ctx, id, app, nil)
	// We don't want to consider a 404 an error in some cases and thus the delineation.
	// Check if app's ID is set to ensure that app exists
	if err := suppressErrorOn404(resp, err); err != nil {
		return err
	}
	return nil
}

func updateAppByID(ctx context.Context, id string, m interface{}, app okta.App) error {
	_, resp, err := getOktaClientFromMetadata(m).Application.UpdateApplication(ctx, id, app)
	// We don't want to consider a 404 an error in some cases and thus the delineation
	if err := suppressErrorOn404(resp, err); err != nil {
		return err
	}
	return nil
}

func handleAppGroups(ctx context.Context, id string, d *schema.ResourceData, client *okta.Client) []func() error {
	existingGroup, _, _ := client.Application.ListApplicationGroupAssignments(ctx, id, nil)
	var (
		asyncActionList []func() error
		groupIDList     []string
	)

	if arr, ok := d.GetOk("groups"); ok {
		rawArr := arr.(*schema.Set).List()
		groupIDList = make([]string, len(rawArr))

		for i, gID := range rawArr {
			groupID := gID.(string)
			groupIDList[i] = groupID

			if !containsGroup(existingGroup, groupID) {
				asyncActionList = append(asyncActionList, func() error {
					_, resp, err := client.Application.CreateApplicationGroupAssignment(ctx, id,
						groupID, okta.ApplicationGroupAssignment{})
					return responseErr(resp, err)
				})
			}
		}
	}

	for _, group := range existingGroup {
		if !contains(groupIDList, group.Id) {
			groupID := group.Id
			asyncActionList = append(asyncActionList, func() error {
				return suppressErrorOn404(client.Application.DeleteApplicationGroupAssignment(ctx, id, groupID))
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
		if user.Id == id && user.Scope == userScope {
			return true
		}
	}
	return false
}

// Handles the assigning of groups and users to Applications. Does so asynchronously.
func handleAppGroupsAndUsers(ctx context.Context, id string, d *schema.ResourceData, m interface{}) error {
	var wg sync.WaitGroup
	resultChan := make(chan []*result, 1)
	client := getOktaClientFromMetadata(m)

	groupHandlers := handleAppGroups(ctx, id, d, client)
	userHandlers := handleAppUsers(ctx, id, d, client)
	con := getParallelismFromMetadata(m)
	promiseAll(con, &wg, resultChan, append(groupHandlers, userHandlers...)...)
	wg.Wait()

	return getPromiseError(<-resultChan, "failed to associate user or groups with application")
}

func handleAppUsers(ctx context.Context, id string, d *schema.ResourceData, client *okta.Client) []func() error {
	// Looking upstream for existing user's, rather then the config for accuracy.
	existingUsers, _, _ := client.Application.ListApplicationUsers(ctx, id, nil)
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
					_, _, err := client.Application.AssignUserToApplication(ctx, id, okta.AppUser{
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
		if user.Scope == userScope {
			if !contains(userIDList, user.Id) {
				userID := user.Id
				asyncActionList = append(asyncActionList, func() error {
					return suppressErrorOn404(client.Application.DeleteApplicationUser(ctx, id, userID, nil))
				})
			}
		}
	}

	return asyncActionList
}

func setAppStatus(ctx context.Context, d *schema.ResourceData, client *okta.Client, status string) error {
	desiredStatus := d.Get("status").(string)
	if status == desiredStatus {
		return nil
	}
	if desiredStatus == statusInactive {
		return responseErr(client.Application.DeactivateApplication(ctx, d.Id()))
	}
	return responseErr(client.Application.ActivateApplication(ctx, d.Id()))
}

func syncGroupsAndUsers(ctx context.Context, id string, d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	// Temporary high limit to avoid issues short term. Need to support pagination here
	userList, _, err := client.Application.ListApplicationUsers(ctx, id, &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return fmt.Errorf("failed to list application users: %v", err)
	}
	// Temporary high limit to avoid issues short term. Need to support pagination here
	groupList, _, err := client.Application.ListApplicationGroupAssignments(ctx, id, &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return fmt.Errorf("failed to list application group assignments: %v", err)
	}
	flatGroupList := make([]interface{}, len(groupList))

	for i, g := range groupList {
		flatGroupList[i] = g.Id
	}

	var flattenedUserList []interface{}

	for _, user := range userList {
		if user.Scope == userScope {
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
	payload, _ := json.Marshal(flatMap)
	return d.Set("app_settings_json", string(payload))
}

func syncSamlSettings(d *schema.ResourceData, set *okta.SamlApplicationSettings) error {
	_ = d.Set("default_relay_state", set.SignOn.DefaultRelayState)
	_ = d.Set("sso_url", set.SignOn.SsoAcsUrl)
	_ = d.Set("recipient", set.SignOn.Recipient)
	_ = d.Set("destination", set.SignOn.Destination)
	_ = d.Set("audience", set.SignOn.Audience)
	_ = d.Set("idp_issuer", set.SignOn.IdpIssuer)
	_ = d.Set("subject_name_id_template", set.SignOn.SubjectNameIdTemplate)
	_ = d.Set("subject_name_id_format", set.SignOn.SubjectNameIdFormat)
	_ = d.Set("response_signed", set.SignOn.ResponseSigned)
	_ = d.Set("assertion_signed", set.SignOn.AssertionSigned)
	_ = d.Set("signature_algorithm", set.SignOn.SignatureAlgorithm)
	_ = d.Set("digest_algorithm", set.SignOn.DigestAlgorithm)
	_ = d.Set("honor_force_authn", set.SignOn.HonorForceAuthn)
	_ = d.Set("authn_context_class_ref", set.SignOn.AuthnContextClassRef)

	if set.SignOn.AllowMultipleAcsEndpoints != nil {
		if *set.SignOn.AllowMultipleAcsEndpoints {
			acsEndpointsObj := set.SignOn.AcsEndpoints
			if len(acsEndpointsObj) > 0 {
				acsEndpoints := make([]string, len(acsEndpointsObj))
				for i := range acsEndpointsObj {
					acsEndpoints[i] = acsEndpointsObj[i].Url
				}
				_ = d.Set("acs_endpoints", convertStringSetToInterface(acsEndpoints))
			}
		} else {
			_ = d.Set("acs_endpoints", nil)
		}
	}

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

func deleteApplication(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	if d.Get("status").(string) == statusActive {
		_, err := client.Application.DeactivateApplication(ctx, d.Id())
		if err != nil {
			return err
		}
	}
	_, err := client.Application.DeleteApplication(ctx, d.Id())
	return err
}
