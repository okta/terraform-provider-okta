package okta

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

type (
	applicationMap struct {
		Type               string
		RequiredGrantTypes []string
		ValidGrantTypes    []string
	}
)

// I wish the SDK provided these
const (
	authorizationCode string = "authorization_code"
	implicit          string = "implicit"
	password          string = "password"
	refreshToken      string = "refresh_token"
	clientCredentials string = "client_credentials"
)

// Building out structure for the conditional validation logic. It looks like customizing the diff
// is the best way to implement this logic, as it needs to introspect.
// NOTE: opened a ticket to Okta to fix their docs, they are off.
// https://developer.okta.com/docs/api/resources/apps#credentials-settings-details
var appGrantTypeMap = map[string]*applicationMap{
	"web": {
		RequiredGrantTypes: []string{
			authorizationCode,
		},
		ValidGrantTypes: []string{
			authorizationCode,
			implicit,
			refreshToken,
			clientCredentials,
		},
	},
	"native": {
		Type: "native",
		RequiredGrantTypes: []string{
			authorizationCode,
		},
		ValidGrantTypes: []string{
			authorizationCode,
			implicit,
			refreshToken,
			password,
		},
	},
	"browser": {
		ValidGrantTypes: []string{
			implicit,
			authorizationCode,
		},
	},
	"service": {
		ValidGrantTypes: []string{
			clientCredentials,
			implicit,
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
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, v interface{}) error {
			// Force new if omit_secret goes from true to false
			if d.Id() != "" {
				old, new := d.GetChange("omit_secret")
				if old.(bool) && !new.(bool) {
					return d.ForceNew("omit_secret")
				}
			}
			return nil
		},
		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Schema: buildAppSchema(map[string]*schema.Schema{
			"type": {
				Type:             schema.TypeString,
				ValidateDiagFunc: stringInSlice([]string{"web", "native", "browser", "service"}),
				Required:         true,
				ForceNew:         true,
				Description:      "The type of client application.",
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
			"custom_client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"client_id"},
				Description:   "**Deprecated** This property allows you to set your client_id during creation. NOTE: updating after creation will be a no-op, use client_id for that behavior instead.",
				Deprecated:    "This field is being replaced by client_id. Please set that field instead.",
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
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{"none", "client_secret_post", "client_secret_basic", "client_secret_jwt", "private_key_jwt"}),
				Default:          "client_secret_basic",
				Description:      "Requested authentication method for the token endpoint.",
			},
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
			"redirect_uris": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of URIs for use in the redirect-based flow. This is required for all application types except service. Note: see okta_app_oauth_redirect_uri for appending to this list in a decentralized way.",
			},
			"post_logout_redirect_uris": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of URIs for redirection after logout",
			},
			"response_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: stringInSlice([]string{"code", "token", "id_token"}),
				},
				Optional:    true,
				Description: "List of OAuth 2.0 response type strings.",
			},
			"grant_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: stringInSlice([]string{authorizationCode, implicit, password, refreshToken, clientCredentials}),
				},
				Optional:    true,
				Description: "List of OAuth 2.0 grant types. Conditional validation params found here https://developer.okta.com/docs/api/resources/apps#credentials-settings-details. Defaults to minimum requirements per app type.",
			},
			// "Early access" properties.. looks to be in beta which requires opt-in per account
			"tos_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "*Early Access Property*. URI to web page providing client tos (terms of service).",
			},
			"policy_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "*Early Access Property*. URI to web page providing client policy document.",
			},
			"consent_method": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "TRUSTED",
				ValidateDiagFunc: stringInSlice([]string{"REQUIRED", "TRUSTED"}),
				Description:      "*Early Access Property*. Indicates whether user consent is required or implicit. Valid values: REQUIRED, TRUSTED. Default value is TRUSTED",
			},
			"issuer_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{"CUSTOM_URL", "ORG_URL"}),
				Default:          "ORG_URL",
				Description:      "*Early Access Property*. Indicates whether the Okta Authorization Server uses the original Okta org domain URL or a custom domain URL as the issuer of ID token for this client.",
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
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Key type",
							ValidateDiagFunc: stringInSlice([]string{"RSA"}),
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
		}),
	}
}

func resourceAppOAuthCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	if err := validateGrantTypes(d); err != nil {
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
		// Needs to be set immediately, not provided again after this
		_ = d.Set("client_secret", app.Credentials.OauthClient.ClientSecret)
	}
	err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to handle groups and users for oauth application: %v", err)
	}
	return resourceAppOAuthRead(ctx, d, m)
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
	var rawProfile string
	if app.Profile != nil {
		p, _ := json.Marshal(app.Profile)
		rawProfile = string(p)
	}
	_ = d.Set("name", app.Name)
	_ = d.Set("status", app.Status)
	_ = d.Set("sign_on_mode", app.SignOnMode)
	_ = d.Set("label", app.Label)
	_ = d.Set("profile", rawProfile)
	_ = d.Set("type", app.Settings.OauthClient.ApplicationType)
	// Not setting client_secret, it is only provided on create for auth methods that require it
	_ = d.Set("client_id", app.Credentials.OauthClient.ClientId)
	_ = d.Set("token_endpoint_auth_method", app.Credentials.OauthClient.TokenEndpointAuthMethod)
	_ = d.Set("auto_key_rotation", app.Credentials.OauthClient.AutoKeyRotation)
	_ = d.Set("client_uri", app.Settings.OauthClient.ClientUri)
	_ = d.Set("logo_uri", app.Settings.OauthClient.LogoUri)
	_ = d.Set("tos_uri", app.Settings.OauthClient.TosUri)
	_ = d.Set("policy_uri", app.Settings.OauthClient.PolicyUri)
	_ = d.Set("login_uri", app.Settings.OauthClient.InitiateLoginUri)
	_ = d.Set("auto_submit_toolbar", app.Visibility.AutoSubmitToolbar)
	_ = d.Set("hide_ios", app.Visibility.Hide.IOS)
	_ = d.Set("hide_web", app.Visibility.Hide.Web)
	if app.Settings.OauthClient.ConsentMethod != "" { // Early Access Property, might be empty
		_ = d.Set("consent_method", app.Settings.OauthClient.ConsentMethod)
	}
	if app.Settings.OauthClient.IssuerMode != "" {
		_ = d.Set("issuer_mode", app.Settings.OauthClient.IssuerMode)
	}

	// If this is ever changed omit it.
	if d.Get("omit_secret").(bool) {
		_ = d.Set("client_secret", "")
	}

	if app.Settings.OauthClient.JWKS != nil {
		jwks := app.Settings.OauthClient.JWKS.Keys
		arr := make([]map[string]interface{}, len(jwks))
		for i, jwk := range jwks {
			arr[i] = map[string]interface{}{
				"kty": jwk.Type,
				"kid": jwk.ID,
				"e":   jwk.Exponent,
				"n":   jwk.Modulus,
			}
		}
		err = setNonPrimitives(d, map[string]interface{}{"jwks": arr})
		if err != nil {
			return diag.Errorf("failed to set OAuth application properties: %v", err)
		}
	}

	respTypes := make([]string, len(app.Settings.OauthClient.ResponseTypes))
	for i := range app.Settings.OauthClient.ResponseTypes {
		respTypes[i] = string(*app.Settings.OauthClient.ResponseTypes[i])
	}
	grantTypes := make([]string, len(app.Settings.OauthClient.GrantTypes))
	for i := range app.Settings.OauthClient.GrantTypes {
		grantTypes[i] = string(*app.Settings.OauthClient.GrantTypes[i])
	}
	err = syncGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to sync groups and users for OAuth application: %v", err)
	}
	aggMap := map[string]interface{}{
		"redirect_uris":             convertStringSetToInterface(app.Settings.OauthClient.RedirectUris),
		"response_types":            convertStringSetToInterface(respTypes),
		"grant_types":               convertStringSetToInterface(grantTypes),
		"post_logout_redirect_uris": convertStringSetToInterface(app.Settings.OauthClient.PostLogoutRedirectUris),
	}
	err = setNonPrimitives(d, aggMap)
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
	app := buildAppOAuth(d)
	_, _, err := client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update OAuth application: %v", err)
	}
	err = setAppStatus(ctx, d, client, app.Status)
	if err != nil {
		return diag.Errorf("failed to set OAuth application status: %v", err)
	}
	err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to handle groups and users for OAuth application: %v", err)
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

	// Need to a bool pointer, it appears the Okta SDK uses this as a way to avoid false being omitted.
	keyRotation := d.Get("auto_key_rotation").(bool)
	appType := d.Get("type").(string)
	grantTypes := convertInterfaceToStringSet(d.Get("grant_types"))
	responseTypes := convertInterfaceToStringSetNullable(d.Get("response_types"))

	// If grant_types are not set, we default to the bare minimum.
	if len(grantTypes) < 1 {
		appMap := appGrantTypeMap[appType]

		if appMap.RequiredGrantTypes == nil {
			grantTypes = appMap.ValidGrantTypes
		} else {
			grantTypes = appMap.RequiredGrantTypes
		}
	}

	// Letting users override response types as well but we properly default them when missing.
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
	app.Credentials = &okta.OAuthApplicationCredentials{
		OauthClient: &okta.ApplicationCredentialsOAuthClient{
			AutoKeyRotation:         &keyRotation,
			ClientId:                d.Get("client_id").(string),
			TokenEndpointAuthMethod: authMethod,
		},
	}

	if sec, ok := d.GetOk("client_basic_secret"); ok {
		app.Credentials.OauthClient.ClientSecret = sec.(string)
	}

	if cid, ok := d.GetOk("custom_client_id"); ok {
		app.Credentials.OauthClient.ClientId = cid.(string)
	}

	oktaRespTypes := make([]*okta.OAuthResponseType, len(responseTypes))
	for i := range responseTypes {
		rt := okta.OAuthResponseType(responseTypes[i])
		oktaRespTypes[i] = &rt
	}
	oktaGrantTypes := make([]*okta.OAuthGrantType, len(grantTypes))
	for i := range grantTypes {
		gt := okta.OAuthGrantType(grantTypes[i])
		oktaGrantTypes[i] = &gt
	}

	app.Settings = &sdk.OpenIdConnectApplicationSettings{
		OauthClient: &sdk.OpenIdConnectApplicationSettingsClient{
			OpenIdConnectApplicationSettingsClient: okta.OpenIdConnectApplicationSettingsClient{
				ApplicationType:        appType,
				ClientUri:              d.Get("client_uri").(string),
				ConsentMethod:          d.Get("consent_method").(string),
				GrantTypes:             oktaGrantTypes,
				InitiateLoginUri:       d.Get("login_uri").(string),
				LogoUri:                d.Get("logo_uri").(string),
				PolicyUri:              d.Get("policy_uri").(string),
				RedirectUris:           convertInterfaceToStringSetNullable(d.Get("redirect_uris")),
				PostLogoutRedirectUris: convertInterfaceToStringSetNullable(d.Get("post_logout_redirect_uris")),
				ResponseTypes:          oktaRespTypes,
				TosUri:                 d.Get("tos_uri").(string),
				IssuerMode:             d.Get("issuer_mode").(string),
			},
		},
	}
	jwks := d.Get("jwks").([]interface{})
	if len(jwks) > 0 {
		keys := make([]*sdk.JWK, len(jwks))
		for i := range jwks {
			keys[i] = &sdk.JWK{
				ID:       d.Get(fmt.Sprintf("jwks.%d.kid", i)).(string),
				Type:     d.Get(fmt.Sprintf("jwks.%d.kty", i)).(string),
				Exponent: d.Get(fmt.Sprintf("jwks.%d.e", i)).(string),
				Modulus:  d.Get(fmt.Sprintf("jwks.%d.n", i)).(string),
			}
		}
		app.Settings.OauthClient.JWKS = &sdk.JWKS{Keys: keys}
	}

	app.Visibility = buildVisibility(d)

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
	appMap := appGrantTypeMap[appType]

	// There is some conditional validation around grant types depending on application type.
	return conditionalValidator("grant_types", appType, appMap.RequiredGrantTypes, appMap.ValidGrantTypes, grantTypeList)
}
