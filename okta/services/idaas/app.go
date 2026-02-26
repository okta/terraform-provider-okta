package idaas

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

var (
	baseAppSchema = map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the app.",
		},
		"label": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Application's display name.",
		},
		"sign_on_mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Sign on mode of application.",
		},
		"status": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     StatusActive,
			Description: "Status of application. By default, it is `ACTIVE`",
		},
		"logo": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: utils.LogoFileIsValid(),
			Description:      "Local file path to the logo. The file must be in PNG, JPG, or GIF format, and less than 1 MB in size.",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return new == ""
			},
			StateFunc: utils.LocalFileStateFunc,
		},
		"logo_url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "URL of the application's logo",
		},
		"admin_note": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Application notes for admins.",
			ValidateDiagFunc: utils.StrMaxLength(250),
		},
		"enduser_note": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Application notes for end users.",
			ValidateDiagFunc: utils.StrMaxLength(250),
		},
		"app_links_json": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Displays specific appLinks for the app. The value for each application link should be boolean.",
			ValidateDiagFunc: stringIsJSON,
			StateFunc:        utils.NormalizeDataJSON,
			DiffSuppressFunc: utils.NoChangeInObjectFromUnmarshaledJSON,
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
			Description: "Enable self service. Default is `false`",
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
			Deprecated:  "Because users has been removed, this attribute is a no op and will be removed",
		},
		"skip_groups": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Ignore groups sync. This is a temporary solution until 'groups' field is supported in all the app-like resources",
			Default:     false,
			Deprecated:  "Because groups has been removed, this attribute is a no op and will be removed",
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
			Description: "Username template. Default: `${source.login}`",
		},
		"user_name_template_suffix": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Username template suffix",
		},
		"user_name_template_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "BUILT_IN",
			Description: "Username template type. Default: `BUILT_IN`",
		},
		"user_name_template_push_status": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Push username on update. Valid values: `PUSH` and `DONT_PUSH`",
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
		// lintignore:R001
		_ = d.Set(v, true)
	}
	return []*schema.ResourceData{d}, nil
}

func appRead(d *schema.ResourceData, name, status, signOn, label string, accy *sdk.ApplicationAccessibility,
	vis *sdk.ApplicationVisibility, notes *sdk.ApplicationSettingsNotes,
) {
	_ = d.Set("name", name)
	_ = d.Set("status", status)
	_ = d.Set("sign_on_mode", signOn)
	_ = d.Set("label", label)
	if accy != nil {
		_ = d.Set("accessibility_self_service", accy.SelfService)
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

func BuildAppSchema(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return utils.BuildSchema(baseAppSchema, baseAppSwaSchema, appSchema)
}

func BuildAppSchemaWithVisibility(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return utils.BuildSchema(baseAppSchema, appVisibilitySchema, appSchema)
}

func BuildAppSwaSchema(appSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return utils.BuildSchema(baseAppSchema, appVisibilitySchema, baseAppSwaSchema, appSchema)
}

func BuildSchemeAppCreds(d *schema.ResourceData) *sdk.SchemeApplicationCredentials {
	revealPass := d.Get("reveal_password").(bool)
	return &sdk.SchemeApplicationCredentials{
		RevealPassword:   &revealPass,
		Scheme:           d.Get("credentials_scheme").(string),
		UserNameTemplate: BuildUserNameTemplate(d),
		UserName:         d.Get("shared_username").(string),
		Password: &sdk.PasswordCredential{
			Value: d.Get("shared_password").(string),
		},
	}
}

func BuildUserNameTemplate(d *schema.ResourceData) *sdk.ApplicationCredentialsUsernameTemplate {
	return &sdk.ApplicationCredentialsUsernameTemplate{
		Template:   d.Get("user_name_template").(string),
		Type:       d.Get("user_name_template_type").(string),
		Suffix:     d.Get("user_name_template_suffix").(string),
		PushStatus: d.Get("user_name_template_push_status").(string),
	}
}

func BuildAppAccessibility(d *schema.ResourceData) *sdk.ApplicationAccessibility {
	return &sdk.ApplicationAccessibility{
		SelfService:      utils.BoolPtr(d.Get("accessibility_self_service").(bool)),
		ErrorRedirectUrl: d.Get("accessibility_error_redirect_url").(string),
		LoginRedirectUrl: d.Get("accessibility_login_redirect_url").(string),
	}
}

func BuildAppVisibility(d *schema.ResourceData) *sdk.ApplicationVisibility {
	autoSubmit := d.Get("auto_submit_toolbar").(bool)
	hideMobile := d.Get("hide_ios").(bool)
	hideWeb := d.Get("hide_web").(bool)
	appVis := &sdk.ApplicationVisibility{
		AutoSubmitToolbar: &autoSubmit,
		Hide: &sdk.ApplicationVisibilityHide{
			IOS: &hideMobile,
			Web: &hideWeb,
		},
	}
	if appLinks, ok := d.GetOk("app_links_json"); ok {
		_ = json.Unmarshal([]byte(appLinks.(string)), &appVis.AppLinks)
	}
	return appVis
}

func BuildAppNotes(d *schema.ResourceData) *sdk.ApplicationSettingsNotes {
	n := &sdk.ApplicationSettingsNotes{}
	admin, ok := d.GetOk("admin_note")
	if ok {
		n.Admin = utils.StringPtr(admin.(string))
	}
	enduser, ok := d.GetOk("enduser_note")
	if ok {
		n.Enduser = utils.StringPtr(enduser.(string))
	}
	return n
}

func BuildAppSettings(d *schema.ResourceData) *sdk.ApplicationSettingsApplication {
	settings := sdk.ApplicationSettingsApplication(map[string]interface{}{})
	if appSettings, ok := d.GetOk("app_settings_json"); ok {
		payload := map[string]interface{}{}
		_ = json.Unmarshal([]byte(appSettings.(string)), &payload)
		settings = payload
	}
	return &settings
}

func fetchApp(ctx context.Context, d *schema.ResourceData, m interface{}, app sdk.App) error {
	return fetchAppByID(ctx, d.Id(), m, app)
}

func fetchAppByID(ctx context.Context, id string, m interface{}, app sdk.App) error {
	_, resp, err := getOktaClientFromMetadata(m).Application.GetApplication(ctx, id, app, nil)
	// We don't want to consider a 404 an error in some cases and thus the delineation.
	// Check if app's ID is set to ensure that app exists
	return utils.SuppressErrorOn404(resp, err)
}

func updateAppByID(ctx context.Context, id string, m interface{}, app sdk.App) error {
	_, resp, err := getOktaClientFromMetadata(m).Application.UpdateApplication(ctx, id, app)
	// We don't want to consider a 404 an error in some cases and thus the delineation
	return utils.SuppressErrorOn404(resp, err)
}

func listApplicationGroupAssignments(ctx context.Context, client *sdk.Client, id string) ([]*sdk.ApplicationGroupAssignment, *sdk.Response, error) {
	groups, resp, err := client.Application.ListApplicationGroupAssignments(ctx, id, &query.Params{Limit: utils.DefaultPaginationLimit})
	if err != nil {
		return nil, resp, err
	}
	for resp.HasNextPage() {
		var additionalGroups []*sdk.ApplicationGroupAssignment
		resp, err = resp.Next(ctx, &additionalGroups)
		if err != nil {
			return nil, resp, err
		}
		groups = append(groups, additionalGroups...)
	}
	return groups, resp, nil
}

func handleAppLogo(ctx context.Context, d *schema.ResourceData, m interface{}, appID string, _ interface{}) error {
	l, ok := d.GetOk("logo")
	if !ok {
		return nil
	}
	_, err := getOktaClientFromMetadata(m).Application.UploadApplicationLogo(ctx, appID, l.(string))
	return err
}

// setAppSettings available preconfigured SAML and OAuth applications vary wildly on potential app settings, thus
// it is a generic map. This logic simply weeds out any empty string values.
func setAppSettings(d *schema.ResourceData, settings *sdk.ApplicationSettingsApplication) error {
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

func setSamlSettings(d *schema.ResourceData, signOn *sdk.SamlApplicationSettingsSignOn) error {
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
	_ = d.Set("saml_signed_request_enabled", signOn.SamlSignedRequestEnabled)
	if signOn.AllowMultipleAcsEndpoints != nil {
		if *signOn.AllowMultipleAcsEndpoints {
			acsEndpointsObj := signOn.AcsEndpoints
			if len(acsEndpointsObj) > 0 {
				indexSequential := isACSEndpointSequential(acsEndpointsObj)
				if indexSequential {
					acsEndpoints := make([]string, len(acsEndpointsObj))
					for i := range acsEndpointsObj {
						acsEndpoints[i] = acsEndpointsObj[i].Url
					}
					_ = d.Set("acs_endpoints", acsEndpoints)
					_ = d.Set("acs_endpoints_indices", nil)
				} else {
					acsList := make([]map[string]interface{}, 0, len(acsEndpointsObj))
					for _, endpoint := range acsEndpointsObj {
						acsList = append(acsList, map[string]interface{}{
							"url":   endpoint.Url,
							"index": *endpoint.IndexPtr,
						})
					}

					sort.Slice(acsList, func(i, j int) bool {
						return acsList[i]["index"].(int64) < acsList[j]["index"].(int64)
					})
					_ = d.Set("acs_endpoints_indices", acsList)
					_ = d.Set("acs_endpoints", nil)
				}
			}
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
	return utils.SetNonPrimitives(d, map[string]interface{}{
		"attribute_statements": arr,
	})
}

func isACSEndpointSequential(acsEndpointsObj []*sdk.AcsEndpoint) bool {
	indexSequential := true
	for i := range acsEndpointsObj {
		if acsEndpointsObj[i].IndexPtr == nil || *acsEndpointsObj[i].IndexPtr != int64(i) {
			indexSequential = false
			break
		}
	}
	return indexSequential
}

func deleteApplication(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	if d.Get("status").(string) == StatusActive {
		// Okta Core can have eventual consistency issues, use backoff for deactivation as well
		boc := utils.NewExponentialBackOffWithContext(ctx, 20*time.Second)
		err := backoff.Retry(func() error {
			_, err := client.Application.DeactivateApplication(ctx, d.Id())
			if doNotRetry(m, err) {
				return backoff.Permanent(err)
			}
			return err
		}, boc)
		if err != nil {
			return err
		}
	}

	// Okta Core can have eventual consistency issues when deactivating an app
	// which is required before deleting the app.
	boc := utils.NewExponentialBackOffWithContext(ctx, 30*time.Second)
	err := backoff.Retry(func() error {
		_, err := client.Application.DeleteApplication(ctx, d.Id())
		if doNotRetry(m, err) {
			return backoff.Permanent(err)
		}
		return err
	}, boc)

	return err
}

// fetchAppKeys returns the keys from `/api/v1/apps/${applicationId}/credentials/keys` for a given app. Not all fields on the JsonWebKey
// will be set, please consult the documentation (https://developer.okta.com/docs/reference/api/apps/#list-key-credentials-for-application)
// for more information.
func fetchAppKeys(ctx context.Context, m interface{}, appID string) ([]*sdk.JsonWebKey, error) {
	keys, _, err := getOktaClientFromMetadata(m).Application.ListApplicationKeys(ctx, appID)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

// setAppKeys sets the JWKs return by fetchAppKeys on the given resource.
func setAppKeys(d *schema.ResourceData, keys []*sdk.JsonWebKey) error {
	arr := make([]map[string]interface{}, len(keys))

	for i, key := range keys {
		arr[i] = map[string]interface{}{
			"kid":          key.Kid,
			"kty":          key.Kty,
			"use":          key.Use,
			"created":      key.Created.String(),
			"last_updated": key.LastUpdated.String(),
			"expires_at":   key.ExpiresAt.String(),
			"e":            key.E,
			"n":            key.N,
			"x5c":          key.X5c,
			"x5t_s256":     key.X5tS256,
		}
	}

	return d.Set("keys", arr)
}

// AppUpdateStatus will activate/deactivate an app based on a change in the
// status argument. As a convenience it will signal if the caller should also
// make other application updates.
func AppUpdateStatus(ctx context.Context, d *schema.ResourceData, m interface{}) (otherChanges bool, err error) {
	status := d.Get("status").(string)
	statusChanged := d.HasChange("status")
	otherChanges = d.HasChangeExcept("status")
	active := (status == StatusActive)
	inactive := (status == StatusInactive)

	if inactive && statusChanged {
		client := getOktaClientFromMetadata(m)
		err = utils.ResponseErr(client.Application.DeactivateApplication(ctx, d.Id()))
	}

	if active && statusChanged {
		client := getOktaClientFromMetadata(m)
		err = utils.ResponseErr(client.Application.ActivateApplication(ctx, d.Id()))
	}

	return otherChanges, err
}
