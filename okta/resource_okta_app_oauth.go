package okta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

type (
	applicationMap struct {
		RequiredGrantTypes []string
		ValidGrantTypes    []string
	}
)

const (
	authorizationCode = "authorization_code"
	interactionCode   = "interaction_code"
	implicit          = "implicit"
	password          = "password"
	refreshToken      = "refresh_token"
	clientCredentials = "client_credentials"
	tokenExchange     = "urn:ietf:params:oauth:grant-type:token-exchange"
	saml2Bearer       = "urn:ietf:params:oauth:grant-type:saml2-bearer"
	deviceCode        = "urn:ietf:params:oauth:grant-type:device_code"
)

// Building out structure for the conditional validation logic. It looks like customizing the diff
// is the best way to implement this logic, as it needs to introspect.
// NOTE: opened a ticket to Okta to fix their docs, they are off.
// https://developer.okta.com/docs/api/resources/apps#credentials-settings-details
var appRequirementsByType = map[string]*applicationMap{
	"web": {
		RequiredGrantTypes: []string{
			authorizationCode,
		},
		ValidGrantTypes: []string{
			authorizationCode,
			implicit,
			refreshToken,
			clientCredentials,
			saml2Bearer,
			tokenExchange,
			deviceCode,
			interactionCode,
		},
	},
	"native": {
		RequiredGrantTypes: []string{
			authorizationCode,
		},
		ValidGrantTypes: []string{
			authorizationCode,
			implicit,
			refreshToken,
			password,
			saml2Bearer,
			tokenExchange,
			deviceCode,
			interactionCode,
		},
	},
	"browser": {
		ValidGrantTypes: []string{
			implicit,
			authorizationCode,
			refreshToken,
			saml2Bearer,
			tokenExchange,
			deviceCode,
			interactionCode,
		},
	},
	"service": {
		ValidGrantTypes: []string{
			clientCredentials,
			implicit,
			saml2Bearer,
			tokenExchange,
			deviceCode,
			interactionCode,
		},
		RequiredGrantTypes: []string{
			clientCredentials,
		},
	},
}

func resourceAppOAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppOAuthCreate,
		ReadContext:   resourceAppOAuthRead,
		UpdateContext: resourceAppOAuthUpdate,
		DeleteContext: resourceAppOAuthDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
		},
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, v interface{}) error {
			// Force new if omit_secret goes from true to false
			if d.Id() != "" {
				oldValue, newValue := d.GetChange("omit_secret")
				if oldValue.(bool) && !newValue.(bool) {
					return d.ForceNew("omit_secret")
				}
			}
			return nil
		},
		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Schema: buildAppSchema(map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of client application.",
			},
			"client_id": {
				Type: schema.TypeString,
				// This field is Optional + Computed because okta automatically sets the
				// client_id value if none is specified during creation.
				// If the client_id is set after creation, the resource will be recreated only if its different from
				// the computed client_id.
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "OAuth client ID. If set during creation, app is created with this id.",
			},
			"omit_secret": {
				Type:     schema.TypeBool,
				Optional: true,
				// No ForceNew to avoid recreating when going from false => true
				Description: "This tells the provider not to persist the application's secret to state. If this is ever changes from true => false your app will be recreated.",
				Default:     false,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "OAuth client secret key. This will be in plain text in your statefile unless you set omit_secret above.",
			},
			"client_basic_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "OAuth client secret key, this can be set when token_endpoint_auth_method is client_secret_basic.",
			},
			"token_endpoint_auth_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "client_secret_basic",
				Description: "Requested authentication method for the token endpoint.",
			},
			// API docs say that auto_key_rotation will alwas be set true if it
			// is missing on input therefore we can declare it's default to be
			// true in the schema.
			"auto_key_rotation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Requested key rotation mode.",
			},
			"client_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URI to a web page providing information about the client.",
			},
			"logo_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URI that references a logo for the client.",
			},
			"login_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URI that initiates login.",
			},
			"login_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of Idp-Initiated login that the client supports, if any",
				Default:     "DISABLED",
			},
			"login_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of scopes to use for the request",
			},
			"pkce_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Require Proof Key for Code Exchange (PKCE) for additional verification key rotation mode. See: https://developer.okta.com/docs/reference/api/apps/#oauth-credential-object",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// when pkce_required is not set in the HCL
					if old == "true" && new == "false" {
						return true
					}
					return false
				},
			},
			"redirect_uris": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of URIs for use in the redirect-based flow. This is required for all application types except service. Note: see okta_app_oauth_redirect_uri for appending to this list in a decentralized way.",
			},
			"wildcard_redirect": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "*Early Access Property*. Indicates if the client is allowed to use wildcard matching of redirect_uris",
				Default:     "DISABLED",
			},
			"post_logout_redirect_uris": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of URIs for redirection after logout. Note: see okta_app_oauth_post_logout_redirect_uri for appending to this list in a decentralized way.",
			},
			"response_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of OAuth 2.0 response type strings.",
			},
			"grant_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of OAuth 2.0 grant types. Conditional validation params found here https://developer.okta.com/docs/api/resources/apps#credentials-settings-details. Defaults to minimum requirements per app type.",
			},
			"tos_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URI to web page providing client tos (terms of service).",
			},
			"policy_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URI to web page providing client policy document.",
			},
			"consent_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "TRUSTED",
				Description: "*Early Access Property*. Indicates whether user consent is required or implicit. Valid values: REQUIRED, TRUSTED. Default value is TRUSTED",
			},
			"issuer_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ORG_URL",
				Description: "*Early Access Property*. Indicates whether the Okta Authorization Server uses the original Okta org domain URL or a custom domain URL as the issuer of ID token for this client.",
			},
			"refresh_token_rotation": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "*Early Access Property* Refresh token rotation behavior",
			},
			"refresh_token_leeway": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "*Early Access Property* Grace period for token rotation",
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
				Default:     true,
				Description: "Do not display application icon on mobile app",
			},
			"hide_web": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Do not display application icon to users",
			},
			"profile": {
				Type:             schema.TypeString,
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        normalizeDataJSON,
				Optional:         true,
				Description:      "Custom JSON that represents an OAuth application's profile",
			},
			"jwks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kid": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key ID",
						},
						"kty": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key type",
						},
						"e": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "RSA Exponent",
						},
						"n": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "RSA Modulus",
						},
					},
				},
			},
			"implicit_assignment": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "*Early Access Property*. Enable Federation Broker Mode.",
			},
			"groups_claim": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Description: "Groups claim for an OpenID Connect client application",
				Optional:    true,
				Elem:        groupsClaimResource,
			},
			"app_settings_json": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Application settings in JSON format",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        normalizeDataJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"authentication_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of this apps authentication policy",
			},
			"jwks_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL reference to JWKS",
			},
		}),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Read:   schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

var groupsClaimResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"type": {
			Description: "Groups claim type.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"filter_type": {
			Description: "Groups claim filter. Can only be set if type is FILTER.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"name": {
			Description: "Name of the claim that will be used in the token.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"value": {
			Description: "Value of the claim. Can be an Okta Expression Language statement that evaluates at the time the token is minted.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"issuer_mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Issuer mode inherited from OAuth App",
		},
	},
}

func resourceAppOAuthCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	if err := validateGrantTypes(d); err != nil {
		return diag.Errorf("failed to create OAuth application: %v", err)
	}
	if err := validateAppOAuth(d); err != nil {
		return diag.Errorf("failed to create OAuth application: %v", err)
	}
	app := buildAppOAuth(d)
	activate := d.Get("status").(string) == statusActive
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create OAuth application: %v", err)
	}
	d.SetId(app.Id)
	if !d.Get("omit_secret").(bool) {
		_ = d.Set("client_secret", app.Credentials.OauthClient.ClientSecret)
	}
	err = handleAppLogo(ctx, d, m, app.Id, app.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for OAuth application: %v", err)
	}
	err = setAppOauthGroupsClaim(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to update groups claim for an OAuth application: %v", err)
	}
	err = createOrUpdateAuthenticationPolicy(ctx, d, m, app.Id)
	if err != nil {
		return diag.Errorf("failed to set authentication policy for an OAuth application: %v", err)
	}
	return resourceAppOAuthRead(ctx, d, m)
}

func setAppOauthGroupsClaim(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	raw, ok := d.GetOk("groups_claim")
	if !ok {
		return nil
	}
	groupsClaim := raw.(*schema.Set).List()[0].(map[string]interface{})
	gc := &sdk.AppOauthGroupClaim{
		Name:  groupsClaim["name"].(string),
		Value: groupsClaim["value"].(string),
	}
	if d.Get("issuer_mode") != nil {
		gc.IssuerMode = d.Get("issuer_mode").(string)
	}
	gct := groupsClaim["type"].(string)
	if gct == "FILTER" {
		gc.ValueType = "GROUPS"
		gc.GroupFilterType = groupsClaim["filter_type"].(string)
	} else {
		gc.ValueType = gct
	}
	_, err := getAPISupplementFromMetadata(m).UpdateAppOauthGroupsClaim(ctx, d.Id(), gc)
	return err
}

func updateAppOauthGroupsClaim(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	raw, ok := d.GetOk("groups_claim")
	if !ok {
		return nil
	}
	if len(raw.(*schema.Set).List()) == 0 {
		gc := &sdk.AppOauthGroupClaim{
			IssuerMode: d.Get("issuer_mode").(string),
		}
		_, err := getAPISupplementFromMetadata(m).UpdateAppOauthGroupsClaim(ctx, d.Id(), gc)
		return err
	}
	return setAppOauthGroupsClaim(ctx, d, m)
}

func resourceAppOAuthRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := sdk.NewOpenIdConnectApplication()
	err := fetchApp(ctx, d, m, app)
	if err != nil {
		return diag.Errorf("failed to get OAuth application: %v", err)
	}
	if app.Id == "" {
		d.SetId("")
		return nil
	}
	setAuthenticationPolicy(d, app.Links)
	var rawProfile string
	if app.Profile != nil {
		p, _ := json.Marshal(app.Profile)
		rawProfile = string(p)
	}
	if app.Credentials.UserNameTemplate != nil {
		_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
		_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
		_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
		_ = d.Set("user_name_template_push_status", app.Credentials.UserNameTemplate.PushStatus)
	}
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	_ = d.Set("profile", rawProfile)
	// Not setting client_secret, it is only provided on create and update for auth methods that require it
	if app.Credentials.OauthClient != nil {
		_ = d.Set("client_id", app.Credentials.OauthClient.ClientId)
		_ = d.Set("token_endpoint_auth_method", app.Credentials.OauthClient.TokenEndpointAuthMethod)
		_ = d.Set("auto_key_rotation", app.Credentials.OauthClient.AutoKeyRotation)
		_ = d.Set("pkce_required", app.Credentials.OauthClient.PkceRequired)
	}
	err = setAppSettings(d, app.Settings.App)
	if err != nil {
		return diag.Errorf("failed to set OAuth application settings: %v", err)
	}
	_ = d.Set("logo_url", linksValue(app.Links, "logo", "href"))
	if app.Settings.ImplicitAssignment != nil {
		_ = d.Set("implicit_assignment", *app.Settings.ImplicitAssignment)
	} else {
		_ = d.Set("implicit_assignment", false)
	}
	// If this is ever changed omit it.
	if d.Get("omit_secret").(bool) {
		_ = d.Set("client_secret", "")
	}

	gc, err := flattenGroupsClaim(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("groups_claim", gc)

	return setOAuthClientSettings(d, app.Settings.OauthClient)
}

func flattenGroupsClaim(ctx context.Context, d *schema.ResourceData, m interface{}) (*schema.Set, error) {
	gc, resp, err := getAPISupplementFromMetadata(m).GetAppOauthGroupsClaim(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return nil, fmt.Errorf("failed to get groups claim for OAuth application: %w", err)
	}
	if gc == nil || gc.Name == "" {
		return nil, nil
	}
	elem := map[string]interface{}{
		"name":  gc.Name,
		"value": gc.Value,
		"type":  gc.ValueType,
	}
	if gc.ValueType == "GROUPS" {
		elem["type"] = "FILTER"
		elem["filter_type"] = gc.GroupFilterType
		elem["issuer_mode"] = gc.IssuerMode
	}
	return schema.NewSet(schema.HashResource(groupsClaimResource), []interface{}{elem}), nil
}

func setOAuthClientSettings(d *schema.ResourceData, oauthClient *sdk.OpenIdConnectApplicationSettingsClient) diag.Diagnostics {
	if oauthClient == nil {
		return nil
	}
	_ = d.Set("type", oauthClient.ApplicationType)
	_ = d.Set("client_uri", oauthClient.ClientUri)
	_ = d.Set("logo_uri", oauthClient.LogoUri)
	_ = d.Set("tos_uri", oauthClient.TosUri)
	_ = d.Set("policy_uri", oauthClient.PolicyUri)
	_ = d.Set("login_uri", oauthClient.InitiateLoginUri)
	_ = d.Set("jwks_uri", oauthClient.JwksUri)
	if oauthClient.WildcardRedirect != "" {
		_ = d.Set("wildcard_redirect", oauthClient.WildcardRedirect)
	}
	if oauthClient.ConsentMethod != "" { // Early Access Property, might be empty
		_ = d.Set("consent_method", oauthClient.ConsentMethod)
	}
	if oauthClient.IssuerMode != "" {
		_ = d.Set("issuer_mode", oauthClient.IssuerMode)
	}
	if oauthClient.RefreshToken != nil {
		_ = d.Set("refresh_token_rotation", oauthClient.RefreshToken.RotationType)
		if oauthClient.RefreshToken.LeewayPtr != nil {
			_ = d.Set("refresh_token_leeway", *oauthClient.RefreshToken.LeewayPtr)
		}
	}
	if oauthClient.Jwks != nil {
		jwks := oauthClient.Jwks.Keys
		arr := make([]map[string]interface{}, len(jwks))
		for i, jwk := range jwks {
			arr[i] = map[string]interface{}{
				"kty": jwk.Kty,
				"kid": jwk.Kid,
				"e":   jwk.E,
				"n":   jwk.N,
			}
		}
		err := setNonPrimitives(d, map[string]interface{}{"jwks": arr})
		if err != nil {
			return diag.Errorf("failed to set OAuth application properties: %v", err)
		}
	}

	respTypes := make([]string, len(oauthClient.ResponseTypes))
	for i := range oauthClient.ResponseTypes {
		respTypes[i] = string(*oauthClient.ResponseTypes[i])
	}
	grantTypes := make([]string, len(oauthClient.GrantTypes))
	for i := range oauthClient.GrantTypes {
		grantTypes[i] = string(*oauthClient.GrantTypes[i])
	}
	aggMap := map[string]interface{}{
		"redirect_uris":             oauthClient.RedirectUris,
		"response_types":            convertStringSliceToSet(respTypes),
		"grant_types":               convertStringSliceToSet(grantTypes),
		"post_logout_redirect_uris": convertStringSliceToSet(oauthClient.PostLogoutRedirectUris),
	}
	if oauthClient.IdpInitiatedLogin != nil {
		_ = d.Set("login_mode", oauthClient.IdpInitiatedLogin.Mode)
		aggMap["login_scopes"] = convertStringSliceToSet(oauthClient.IdpInitiatedLogin.DefaultScope)
	}
	err := setNonPrimitives(d, aggMap)
	if err != nil {
		return diag.Errorf("failed to set OAuth application properties: %v", err)
	}
	return nil
}

func resourceAppOAuthUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	if err := validateGrantTypes(d); err != nil {
		return diag.Errorf("failed to update OAuth application: %v", err)
	}
	if err := validateAppOAuth(d); err != nil {
		return diag.Errorf("failed to create OAuth application: %v", err)
	}
	app := buildAppOAuth(d)
	_, _, err := client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update OAuth application: %v", err)
	}
	if !d.Get("omit_secret").(bool) {
		_ = d.Set("client_secret", app.Credentials.OauthClient.ClientSecret)
	}
	err = setAppStatus(ctx, d, client, app.Status)
	if err != nil {
		return diag.Errorf("failed to set OAuth application status: %v", err)
	}
	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, m, app.Id, app.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for OAuth application: %v", err)
		}
	}
	err = updateAppOauthGroupsClaim(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to update groups claim for an OAuth application: %v", err)
	}
	err = createOrUpdateAuthenticationPolicy(ctx, d, m, app.Id)
	if err != nil {
		return diag.Errorf("failed to set authentication policy an OAuth application: %v", err)
	}
	return resourceAppOAuthRead(ctx, d, m)
}

func resourceAppOAuthDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to delete OAuth application: %v", err)
	}
	return nil
}

func buildAppOAuth(d *schema.ResourceData) *sdk.OpenIdConnectApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := sdk.NewOpenIdConnectApplication()
	appType := d.Get("type").(string)
	grantTypes := convertInterfaceToStringSet(d.Get("grant_types"))
	responseTypes := convertInterfaceToStringSetNullable(d.Get("response_types"))

	// If grant_types are not set, we default to the bare minimum.
	if len(grantTypes) < 1 {
		appMap, ok := appRequirementsByType[appType]
		if ok {
			if appMap.RequiredGrantTypes == nil {
				grantTypes = appMap.ValidGrantTypes
			} else {
				grantTypes = appMap.RequiredGrantTypes
			}
		}
	}

	// Letting users override response types as well, but we properly default them when missing.
	if len(responseTypes) < 1 {
		responseTypes = []string{}

		if containsOne(grantTypes, implicit, clientCredentials) {
			responseTypes = append(responseTypes, "token")
		}

		if containsOne(grantTypes, password, authorizationCode, refreshToken) {
			responseTypes = append(responseTypes, "code")
		}
	}

	app.Label = d.Get("label").(string)
	authMethod := d.Get("token_endpoint_auth_method").(string)
	app.Credentials = &sdk.OAuthApplicationCredentials{
		OauthClient: &sdk.ApplicationCredentialsOAuthClient{
			AutoKeyRotation:         boolPtr(d.Get("auto_key_rotation").(bool)),
			ClientId:                d.Get("client_id").(string),
			TokenEndpointAuthMethod: authMethod,
			ClientSecret:            d.Get("client_secret").(string),
		},
		UserNameTemplate: buildUserNameTemplate(d),
	}

	var pkceRequired *bool
	pkceVal := d.GetRawConfig().GetAttr("pkce_required")
	// only explicitly set pkce_required to true for browser and native apps
	// when it isn't set in the HCL
	if pkceVal.IsNull() && (appType == "native" || appType == "browser") {
		pkceRequired = boolPtr(true)
	} else {
		switch {
		case pkceVal.True():
			pkceRequired = boolPtr(true)
		case pkceVal.False():
			pkceRequired = boolPtr(false)
		}
	}
	app.Credentials.OauthClient.PkceRequired = pkceRequired

	if sec, ok := d.GetOk("client_basic_secret"); ok {
		app.Credentials.OauthClient.ClientSecret = sec.(string)
	}

	oktaRespTypes := make([]*sdk.OAuthResponseType, len(responseTypes))
	for i := range responseTypes {
		rt := sdk.OAuthResponseType(responseTypes[i])
		oktaRespTypes[i] = &rt
	}
	oktaGrantTypes := make([]*sdk.OAuthGrantType, len(grantTypes))
	for i := range grantTypes {
		gt := sdk.OAuthGrantType(grantTypes[i])
		oktaGrantTypes[i] = &gt
	}
	app.Settings = &sdk.OpenIdConnectApplicationSettings{
		ImplicitAssignment: boolPtr(d.Get("implicit_assignment").(bool)),
		OauthClient: &sdk.OpenIdConnectApplicationSettingsClient{
			ApplicationType:        appType,
			ClientUri:              d.Get("client_uri").(string),
			ConsentMethod:          d.Get("consent_method").(string),
			GrantTypes:             oktaGrantTypes,
			InitiateLoginUri:       d.Get("login_uri").(string),
			LogoUri:                d.Get("logo_uri").(string),
			PolicyUri:              d.Get("policy_uri").(string),
			RedirectUris:           convertInterfaceToStringArr(d.Get("redirect_uris")),
			PostLogoutRedirectUris: convertInterfaceToStringSetNullable(d.Get("post_logout_redirect_uris")),
			ResponseTypes:          oktaRespTypes,
			TosUri:                 d.Get("tos_uri").(string),
			IssuerMode:             d.Get("issuer_mode").(string),
			IdpInitiatedLogin: &sdk.OpenIdConnectApplicationIdpInitiatedLogin{
				DefaultScope: convertInterfaceToStringSet(d.Get("login_scopes")),
				Mode:         d.Get("login_mode").(string),
			},
			WildcardRedirect: d.Get("wildcard_redirect").(string),
			JwksUri:          d.Get("jwks_uri").(string),
		},
		Notes: buildAppNotes(d),
		App:   buildAppSettings(d),
	}
	jwks := d.Get("jwks").([]interface{})
	if len(jwks) > 0 {
		keys := make([]*sdk.JsonWebKey, len(jwks))
		for i := range jwks {
			keys[i] = &sdk.JsonWebKey{
				Kid: d.Get(fmt.Sprintf("jwks.%d.kid", i)).(string),
				Kty: d.Get(fmt.Sprintf("jwks.%d.kty", i)).(string),
				E:   d.Get(fmt.Sprintf("jwks.%d.e", i)).(string),
				N:   d.Get(fmt.Sprintf("jwks.%d.n", i)).(string),
			}
		}
		app.Settings.OauthClient.Jwks = &sdk.OpenIdConnectApplicationSettingsClientKeys{Keys: keys}
	}

	refresh := &sdk.OpenIdConnectApplicationSettingsRefreshToken{}
	var hasRefresh bool
	for _, grant_type := range grantTypes {
		if grant_type == refreshToken {
			hasRefresh = true
			break
		}
	}
	if rotate, ok := d.GetOk("refresh_token_rotation"); ok {
		refresh.RotationType = rotate.(string)
	}

	leeway, ok := d.GetOk("refresh_token_leeway")
	if ok {
		refresh.LeewayPtr = int64Ptr(leeway.(int))
	} else {
		refresh.LeewayPtr = int64Ptr(0)
	}

	if hasRefresh {
		app.Settings.OauthClient.RefreshToken = refresh
	}

	// TODO: need to put a warning
	// if !hasRefresh && refresh != nil {
	// 	return nil, errors.New("does not have refresh grant type but refresh_token_rotation and refresh_token_leeway exist in payload")
	// }
	// TODO unset refresh_token_rotation, refresh_token_leeway

	app.Visibility = buildAppVisibility(d)
	app.Accessibility = buildAppAccessibility(d)

	if rawAttrs, ok := d.GetOk("profile"); ok {
		var attrs map[string]interface{}
		str := rawAttrs.(string)
		_ = json.Unmarshal([]byte(str), &attrs)
		app.Profile = attrs
	}

	return app
}

func validateGrantTypes(d *schema.ResourceData) error {
	grantTypeList := convertInterfaceToStringSet(d.Get("grant_types"))
	appType := d.Get("type").(string)
	appMap, ok := appRequirementsByType[appType]
	if !ok {
		return nil
	}
	// There is some conditional validation around grant types depending on application type.
	return conditionalValidator("grant_types", appType, appMap.RequiredGrantTypes, appMap.ValidGrantTypes, grantTypeList)
}

func validateAppOAuth(d *schema.ResourceData) error {
	raw, ok := d.GetOk("groups_claim")
	if ok {
		groupsClaim := raw.(*schema.Set).List()[0].(map[string]interface{})
		if groupsClaim["type"].(string) == "EXPRESSION" && groupsClaim["filter_type"].(string) != "" {
			return errors.New("'filter_type' in 'groups_claim' can only be set when 'type' is set to 'FILTER'")
		}
		if groupsClaim["type"].(string) == "FILTER" && groupsClaim["filter_type"].(string) == "" {
			return errors.New("'filter_type' in 'groups_claim' is required when 'type' is set to 'FILTER'")
		}
		if groupsClaim["name"].(string) == "" || groupsClaim["value"].(string) == "" {
			return errors.New("'name' 'value' and in 'groups_claim' should not be empty")
		}
	}
	_, jwks := d.GetOk("jwks")
	_, jwks_uri := d.GetOk("jwks_uri")
	if !(jwks || jwks_uri) && d.Get("token_endpoint_auth_method").(string) == "private_key_jwt" {
		return errors.New("'jwks' or 'jwks_uri' is required when 'token_endpoint_auth_method' is 'private_key_jwt'")
	}
	if d.Get("login_mode").(string) != "DISABLED" {
		if d.Get("login_uri").(string) == "" {
			return errors.New("you have to set up 'login_uri' to configure any 'login_mode' besides 'DISABLED'")
		}
		if d.Get("login_mode").(string) == "OKTA" && len(convertInterfaceToStringSet(d.Get("login_scopes"))) < 1 {
			return errors.New("you have to set up non-empty 'login_scopes' when 'login_mode' is 'OKTA'")
		}
	}
	grantTypes := convertInterfaceToStringSet(d.Get("grant_types"))
	hasImplicit := false
	for _, v := range grantTypes {
		if v == "implicit" {
			hasImplicit = true
			break
		}
	}
	if !hasImplicit {
		return nil
	}
	hasTokenOrTokenID := false
	for _, v := range convertInterfaceToStringSetNullable(d.Get("response_types")) {
		if v == "token" || v == "id_token" {
			hasTokenOrTokenID = true
			break
		}
	}
	if !hasTokenOrTokenID {
		return errors.New("'response_types' must contain at least one of ['token', 'id_token'] when 'grant_types' contains 'implicit'")
	}
	return nil
}
