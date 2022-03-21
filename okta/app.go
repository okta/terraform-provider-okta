package okta

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

var (
	appUserResource = &schema.Resource{
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

	baseAppSchema = map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the app.",
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
			Deprecated:  "The direct configuration of users in this app resource is deprecated, please ensure you use the resource `okta_app_user` for this functionality.",
		},
		"groups": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Groups associated with the application",
			Deprecated:  "The direct configuration of groups in this app resource is deprecated, please ensure you use the resource `okta_app_group_assignments` for this functionality.",
		},
		"status": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          statusActive,
			ValidateDiagFunc: elemInSlice([]string{statusActive, statusInactive}),
			Description:      "Status of application.",
		},
		"logo": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: logoValid(),
			Description:      "Local path to logo of the application.",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return new == ""
			},
			StateFunc: logoStateFunc,
		},
		"logo_url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "URL of the application's logo",
		},
		"admin_note": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Application notes for admins.",
		},
		"enduser_note": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Application notes for end users.",
		},
		"app_links_json": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Displays specific appLinks for the app",
			ValidateDiagFunc: stringIsJSON,
			StateFunc:        normalizeDataJSON,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return new == ""
			},
		},
		"accessibility_login_redirect_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Custom login page URL",
		},
		"accessibility_self_service": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable self service",
		},
		"accessibility_error_redirect_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Custom error page URL",
		},
	}

	skipUsersAndGroupsSchema = map[string]*schema.Schema{
		"skip_users": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Ignore users sync. This is a temporary solution until 'users' field is supported in all the app-like resources",
			Default:     false,
		},
		"skip_groups": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Ignore groups sync. This is a temporary solution until 'groups' field is supported in all the app-like resources",
			Default:     false,
		},
	}

	appVisibilitySchema = map[string]*schema.Schema{
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

	baseAppSwaSchema = map[string]*schema.Schema{
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
			ValidateDiagFunc: elemInSlice([]string{"NONE", "CUSTOM", "BUILT_IN"}),
		},
		"user_name_template_push_status": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Push username on update",
			ValidateDiagFunc: elemInSlice([]string{"DONT_PUSH", "PUSH", ""}),
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return new == ""
			},
		},
	}

	appSamlDiffSuppressFunc = func(k, old, new string, d *schema.ResourceData) bool {
		// Conditional default
		return new == "" && old == "http://www.okta.com/${org.externalKey}"
	}
)

func appImporter(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := strings.Split(d.Id(), "/")
	if len(importID) == 1 {
		return []*schema.ResourceData{d}, nil
	}
	if len(importID) > 3 {
		return nil, errors.New("invalid format used for import ID, format must be 'app_id' or 'app_id/skip_users' or 'app_id/skip_users/skip_groups'")
	}
	d.SetId(importID[0])
	for _, v := range importID[1:] {
		if !isValidSkipArg(v) {
			return nil, fmt.Errorf("'%s' is invalid value to be used as part of import ID, it must be either 'skip_users' or 'skip_groups'", v)
		}
		_ = d.Set(v, true)
	}
	return []*schema.ResourceData{d}, nil
}

func appRead(d *schema.ResourceData, name, status, signOn, label string, accy *okta.ApplicationAccessibility,
	vis *okta.ApplicationVisibility, notes *okta.ApplicationSettingsNotes) {
	_ = d.Set("name", name)
	_ = d.Set("status", status)
	_ = d.Set("sign_on_mode", signOn)
	_ = d.Set("label", label)
	if accy != nil {
		_ = d.Set("accessibility_self_service", *accy.SelfService)
		_ = d.Set("accessibility_error_redirect_url", accy.ErrorRedirectUrl)
		_ = d.Set("accessibility_login_redirect_url", accy.LoginRedirectUrl)
	}
	_ = d.Set("auto_submit_toolbar", vis.AutoSubmitToolbar)
	_ = d.Set("hide_ios", vis.Hide.IOS)
	_ = d.Set("hide_web", vis.Hide.Web)
	if notes != nil {
		_ = d.Set("admin_note", notes.Admin)
		_ = d.Set("enduser_note", notes.Enduser)
	}
	_ = setAppLinks(d, vis.AppLinks)
}

func buildAppSchema(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseAppSchema, skipUsersAndGroupsSchema, baseAppSwaSchema, appSchema)
}

func buildAppSchemaWithVisibility(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseAppSchema, skipUsersAndGroupsSchema, appVisibilitySchema, appSchema)
}

func buildAppSwaSchema(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(baseAppSchema, skipUsersAndGroupsSchema, appVisibilitySchema, baseAppSwaSchema, appSchema)
}

func buildSchemeAppCreds(d *schema.ResourceData) *okta.SchemeApplicationCredentials {
	revealPass := d.Get("reveal_password").(bool)
	return &okta.SchemeApplicationCredentials{
		RevealPassword:   &revealPass,
		Scheme:           d.Get("credentials_scheme").(string),
		UserNameTemplate: buildUserNameTemplate(d),
		UserName:         d.Get("shared_username").(string),
		Password: &okta.PasswordCredential{
			Value: d.Get("shared_password").(string),
		},
	}
}

func buildUserNameTemplate(d *schema.ResourceData) *okta.ApplicationCredentialsUsernameTemplate {
	return &okta.ApplicationCredentialsUsernameTemplate{
		Template:   d.Get("user_name_template").(string),
		Type:       d.Get("user_name_template_type").(string),
		Suffix:     d.Get("user_name_template_suffix").(string),
		PushStatus: d.Get("user_name_template_push_status").(string),
	}
}

func buildAppAccessibility(d *schema.ResourceData) *okta.ApplicationAccessibility {
	return &okta.ApplicationAccessibility{
		SelfService:      boolPtr(d.Get("accessibility_self_service").(bool)),
		ErrorRedirectUrl: d.Get("accessibility_error_redirect_url").(string),
		LoginRedirectUrl: d.Get("accessibility_login_redirect_url").(string),
	}
}

func buildAppVisibility(d *schema.ResourceData) *okta.ApplicationVisibility {
	autoSubmit := d.Get("auto_submit_toolbar").(bool)
	hideMobile := d.Get("hide_ios").(bool)
	hideWeb := d.Get("hide_web").(bool)
	appVis := &okta.ApplicationVisibility{
		AutoSubmitToolbar: &autoSubmit,
		Hide: &okta.ApplicationVisibilityHide{
			IOS: &hideMobile,
			Web: &hideWeb,
		},
	}
	if appLinks, ok := d.GetOk("app_links_json"); ok {
		_ = json.Unmarshal([]byte(appLinks.(string)), &appVis.AppLinks)
	}
	return appVis
}

func buildAppNotes(d *schema.ResourceData) *okta.ApplicationSettingsNotes {
	n := &okta.ApplicationSettingsNotes{}
	admin, ok := d.GetOk("admin_note")
	if ok {
		n.Admin = stringPtr(admin.(string))
	}
	enduser, ok := d.GetOk("enduser_note")
	if ok {
		n.Enduser = stringPtr(enduser.(string))
	}
	return n
}

func buildAppSettings(d *schema.ResourceData) *okta.ApplicationSettingsApplication {
	settings := okta.ApplicationSettingsApplication(map[string]interface{}{})
	if appSettings, ok := d.GetOk("app_settings_json"); ok {
		payload := map[string]interface{}{}
		_ = json.Unmarshal([]byte(appSettings.(string)), &payload)
		settings = payload
	}
	return &settings
}

func fetchApp(ctx context.Context, d *schema.ResourceData, m interface{}, app okta.App) error {
	return fetchAppByID(ctx, d.Id(), m, app)
}

func fetchAppByID(ctx context.Context, id string, m interface{}, app okta.App) error {
	_, resp, err := getOktaClientFromMetadata(m).Application.GetApplication(ctx, id, app, nil)
	// We don't want to consider a 404 an error in some cases and thus the delineation.
	// Check if app's ID is set to ensure that app exists
	return suppressErrorOn404(resp, err)
}

func updateAppByID(ctx context.Context, id string, m interface{}, app okta.App) error {
	_, resp, err := getOktaClientFromMetadata(m).Application.UpdateApplication(ctx, id, app)
	// We don't want to consider a 404 an error in some cases and thus the delineation
	return suppressErrorOn404(resp, err)
}

func handleAppGroups(ctx context.Context, id string, d *schema.ResourceData, client *okta.Client) []func() error {
	if !d.HasChange("groups") {
		return nil
	}
	// temp solution until 'groups' field is supported
	if d.Get("skip_groups").(bool) {
		return nil
	}
	var asyncActionList []func() error

	oldGs, newGs := d.GetChange("groups")
	oldSet := oldGs.(*schema.Set)
	newSet := newGs.(*schema.Set)
	groupsToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	groupsToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

	for i := range groupsToAdd {
		gID := groupsToAdd[i]
		asyncActionList = append(asyncActionList, func() error {
			_, resp, err := client.Application.CreateApplicationGroupAssignment(ctx, id,
				gID, okta.ApplicationGroupAssignment{})
			return responseErr(resp, err)
		})
	}
	for i := range groupsToRemove {
		gID := groupsToRemove[i]
		asyncActionList = append(asyncActionList, func() error {
			return suppressErrorOn404(client.Application.DeleteApplicationGroupAssignment(ctx, id, gID))
		})
	}
	return asyncActionList
}

func listApplicationGroupAssignments(ctx context.Context, client *okta.Client, id string) ([]*okta.ApplicationGroupAssignment, *okta.Response, error) {
	groups, resp, err := client.Application.ListApplicationGroupAssignments(ctx, id, &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return nil, resp, err
	}
	for resp.HasNextPage() {
		var additionalGroups []*okta.ApplicationGroupAssignment
		resp, err = resp.Next(ctx, &additionalGroups)
		if err != nil {
			return nil, resp, err
		}
		groups = append(groups, additionalGroups...)
	}
	return groups, resp, nil
}

func containsAppUser(userList []*okta.AppUser, id string) bool {
	for _, user := range userList {
		if user.Id == id && user.Scope == userScope {
			return true
		}
	}
	return false
}

func shouldUpdateUser(userList []*okta.AppUser, id, username string) bool {
	for _, user := range userList {
		if user.Id == id &&
			user.Scope == userScope &&
			user.Credentials != nil &&
			user.Credentials.UserName != username {
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

func handleAppLogo(ctx context.Context, d *schema.ResourceData, m interface{}, appID string, links interface{}) error {
	l, ok := d.GetOk("logo")
	if !ok {
		return nil
	}
	_, err := getOktaClientFromMetadata(m).Application.UploadApplicationLogo(ctx, appID, l.(string))
	return err
}

func handleAppUsers(ctx context.Context, id string, d *schema.ResourceData, client *okta.Client) []func() error {
	if !d.HasChange("users") {
		return nil
	}
	// temp solution until 'users' field is supported
	if d.Get("skip_users").(bool) {
		return nil
	}
	existingUsers, err := listApplicationUsers(ctx, client, id)
	if err != nil {
		return []func() error{
			func() error { return err },
		}
	}

	var asyncActionList []func() error

	oldUs, newUs := d.GetChange("users")
	oldSet := oldUs.(*schema.Set)
	newSet := newUs.(*schema.Set)
	usersToAdd := newSet.Difference(oldSet).List()
	usersToRemove := oldSet.Difference(newSet).List()

	for i := range usersToAdd {
		userProfile := usersToAdd[i].(map[string]interface{})
		uID := userProfile["id"].(string)
		username := userProfile["username"].(string)
		password := userProfile["password"].(string)
		if shouldUpdateUser(existingUsers, uID, username) {
			asyncActionList = append(asyncActionList, func() error {
				_, _, err := client.Application.UpdateApplicationUser(ctx, id, uID, okta.AppUser{
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
		} else {
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

	for i := range usersToRemove {
		uID := usersToRemove[i].(map[string]interface{})["id"].(string)
		if containsAppUser(existingUsers, uID) {
			asyncActionList = append(asyncActionList, func() error {
				return suppressErrorOn404(client.Application.DeleteApplicationUser(ctx, id, uID, nil))
			})
		}
	}
	return asyncActionList
}

func listApplicationUsers(ctx context.Context, client *okta.Client, id string) ([]*okta.AppUser, error) {
	var resUsers []*okta.AppUser
	users, resp, err := client.Application.ListApplicationUsers(ctx, id, &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return nil, err
	}
	for {
		resUsers = append(resUsers, users...)
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &users)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			break
		}
	}
	return resUsers, nil
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
	ctx = context.WithValue(ctx, retryOnStatusCodes, []int{http.StatusNotFound})

	logger(m).Info("skip values",
		"skip_users", d.Get("skip_users").(bool),
		"skip_groups", d.Get("skip_groups").(bool))

	flatMap := map[string]interface{}{}
	if skipUsers := d.Get("skip_users").(bool); !skipUsers {
		appUsers, err := listApplicationUsers(ctx, getOktaClientFromMetadata(m), id)
		if err != nil {
			return err
		}
		var flattenedUserList []interface{}
		for _, user := range appUsers {
			if user.Scope != userScope {
				continue
			}
			var un, up string
			if user.Credentials != nil {
				getOktaClientFromMetadata(m)
				un = user.Credentials.UserName
				if user.Credentials.Password != nil {
					up = user.Credentials.Password.Value
				}
			}
			flattenedUserList = append(flattenedUserList, map[string]interface{}{
				"id":       user.Id,
				"username": un,
				"scope":    user.Scope,
				"password": up,
			})
		}
		if len(flattenedUserList) > 0 {
			flatMap["users"] = schema.NewSet(schema.HashResource(appUserResource), flattenedUserList)
		}
	}
	if skipGroups := d.Get("skip_groups").(bool); !skipGroups {
		appGroups, _, err := listApplicationGroupAssignments(ctx, getOktaClientFromMetadata(m), id)
		if err != nil {
			return err
		}
		flatGroupList := make([]interface{}, len(appGroups))
		for i := range appGroups {
			flatGroupList[i] = appGroups[i].Id
		}
		if len(flatGroupList) > 0 {
			flatMap["groups"] = schema.NewSet(schema.HashString, flatGroupList)
		}
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

func setAppLinks(d *schema.ResourceData, appLinks map[string]bool) error {
	if len(appLinks) == 0 {
		return nil
	}
	payload, _ := json.Marshal(appLinks)
	return d.Set("app_links_json", string(payload))
}

func setSamlSettings(d *schema.ResourceData, signOn *okta.SamlApplicationSettingsSignOn) error {
	_ = d.Set("default_relay_state", signOn.DefaultRelayState)
	_ = d.Set("sso_url", signOn.SsoAcsUrl)
	_ = d.Set("recipient", signOn.Recipient)
	_ = d.Set("destination", signOn.Destination)
	_ = d.Set("audience", signOn.Audience)
	_ = d.Set("idp_issuer", signOn.IdpIssuer)
	_ = d.Set("subject_name_id_template", signOn.SubjectNameIdTemplate)
	_ = d.Set("subject_name_id_format", signOn.SubjectNameIdFormat)
	_ = d.Set("response_signed", signOn.ResponseSigned)
	_ = d.Set("assertion_signed", signOn.AssertionSigned)
	_ = d.Set("signature_algorithm", signOn.SignatureAlgorithm)
	_ = d.Set("digest_algorithm", signOn.DigestAlgorithm)
	_ = d.Set("honor_force_authn", signOn.HonorForceAuthn)
	_ = d.Set("authn_context_class_ref", signOn.AuthnContextClassRef)
	if signOn.AllowMultipleAcsEndpoints != nil {
		if *signOn.AllowMultipleAcsEndpoints {
			acsEndpointsObj := signOn.AcsEndpoints
			if len(acsEndpointsObj) > 0 {
				acsEndpoints := make([]string, len(acsEndpointsObj))
				for i := range acsEndpointsObj {
					acsEndpoints[i] = acsEndpointsObj[i].Url
				}
				_ = d.Set("acs_endpoints", convertStringSliceToSetNullable(acsEndpoints))
			}
		} else {
			_ = d.Set("acs_endpoints", nil)
		}
	}

	if len(signOn.InlineHooks) > 0 {
		_ = d.Set("inline_hook_id", signOn.InlineHooks[0].Id)
	}
	attrStatements := signOn.AttributeStatements
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
	if signOn.Slo != nil && signOn.Slo.Enabled != nil && *signOn.Slo.Enabled {
		_ = d.Set("single_logout_issuer", signOn.Slo.Issuer)
		_ = d.Set("single_logout_url", signOn.Slo.LogoutUrl)
		if signOn.SpCertificate != nil && len(signOn.SpCertificate.X5c) > 0 {
			_ = d.Set("single_logout_certificate", signOn.SpCertificate.X5c[0])
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

func setAppUsersIDsAndGroupsIDs(ctx context.Context, d *schema.ResourceData, client *okta.Client, id string) error {
	if skipGroups := d.Get("skip_groups").(bool); !skipGroups {
		groups, _, err := listApplicationGroupAssignments(ctx, client, id)
		if err != nil {
			return err
		}
		groupsIDs := make([]string, len(groups))
		for i := range groups {
			groupsIDs[i] = groups[i].Id
		}
		_ = d.Set("groups", convertStringSliceToSet(groupsIDs))
	}
	if skipUsers := d.Get("skip_users").(bool); !skipUsers {
		users, err := listApplicationUsers(ctx, client, id)
		if err != nil {
			return err
		}
		usersIDs := make([]string, len(users))
		for i := range users {
			usersIDs[i] = users[i].Id
		}
		_ = d.Set("users", convertStringSliceToSet(usersIDs))
	}
	return nil
}

func computeFileHash(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	_ = file.Close()
	return hex.EncodeToString(h.Sum(nil))
}

func logoStateFunc(val interface{}) string {
	logoPath := val.(string)
	if logoPath == "" {
		return ""
	}
	return computeFileHash(logoPath)
}
