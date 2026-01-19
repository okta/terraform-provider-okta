package idaas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
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
	jwtBearer         = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	deviceCode        = "urn:ietf:params:oauth:grant-type:device_code"
	oob               = "urn:okta:params:oauth:grant-type:oob"
	otp               = "urn:okta:params:oauth:grant-type:otp"
	mfaOob            = "http://auth0.com/oauth/grant-type/mfa-oob"
	mfaOtp            = "http://auth0.com/oauth/grant-type/mfa-otp"
	ciba              = "urn:openid:params:grant-type:ciba"
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
			oob,
			otp,
			mfaOob,
			mfaOtp,
			ciba,
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
			jwtBearer,
			tokenExchange,
			deviceCode,
			interactionCode,
			oob,
			otp,
			mfaOob,
			mfaOtp,
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
			oob,
			otp,
			mfaOob,
			mfaOtp,
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
			oob,
			otp,
			mfaOob,
			mfaOtp,
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
		Description: `This resource allows you to create and configure an OIDC Application.
-> During an apply if there is change in status the app will first be
activated or deactivated in accordance with the status change. Then, all
other arguments that changed will be applied.`,
		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Schema: BuildAppSchema(map[string]*schema.Schema{
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
				Description: "This tells the provider not manage the client_secret value in state. When this is false (the default), it will cause the auto-generated client_secret to be persisted in the client_secret attribute in state. This also means that every time an update to this app is run, this value is also set on the API. If this changes from false => true, the `client_secret` is dropped from state and the secret at the time of the apply is what remains. If this is ever changes from true => false your app will be recreated, due to the need to regenerate a secret we can store in state.",
				Default:     false,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "OAuth client secret value, this is output only. This will be in plain text in your statefile unless you set omit_secret above.",
			},
			"client_basic_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The user provided OAuth client secret key value, this can be set when token_endpoint_auth_method is client_secret_basic. This does nothing when `omit_secret is set to true.",
			},
			"token_endpoint_auth_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "client_secret_basic",
				Description: "Requested authentication method for the token endpoint, valid values include:  'client_secret_basic', 'client_secret_post', 'client_secret_jwt', 'private_key_jwt', 'none', etc.",
			},
			// API docs say that auto_key_rotation will alwas be set true if it
			// is missing on input therefore we can declare it's default to be
			// true in the schema.
			"auto_key_rotation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: `Requested key rotation mode. If
				auto_key_rotation isn't specified, the client automatically opts in for Okta's
				key rotation. You can update this property via the API or via the administrator
				UI.
				See: https://developer.okta.com/docs/reference/api/apps/#oauth-credential-object"`,
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
				Computed:    true,
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
			"participate_slo": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "*Early Access Property*. Allows the app to participate in front-channel Single Logout. Note: You can only enable participate_slo for web and browser application types. When set to true, frontchannel_logout_uri must also be provided.",
			},
			"frontchannel_logout_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "*Early Access Property*. URL where Okta sends the logout request. Required when participate_slo is true.",
			},
			"frontchannel_logout_session_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "*Early Access Property*. Determines whether Okta sends sid and iss in the logout request.",
			},
			"response_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{ // Normally we don't do input validation, but these values are unlikely to change as they are part of the OAuth 2.0 spec
						"code",
						"token",
						"id_token",
					}, false),
				},
				Optional:    true,
				Computed:    true,
				Description: "List of OAuth 2.0 response type strings. Valid values are any combination of: `code`, `token`, and `id_token`.",
			},
			"grant_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
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
				Default:     "STATIC",
				Description: "*Early Access Property* Refresh token rotation behavior, required with grant types refresh_token",
			},
			"refresh_token_leeway": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "*Early Access Property* Grace period for token rotation, required with grant types refresh_token",
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
				StateFunc:        utils.NormalizeDataJSON,
				Optional:         true,
				Description:      "Custom JSON that represents an OAuth application's profile",
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
			"jwks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "JSON Web Key Set (JWKS) for application. Note: Inline JWKS may have compatibility issues with v6 SDK. Consider using jwks_uri instead.",
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
						"x": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "X coordinate of the elliptic curve point",
						},
						"y": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Y coordinate of the elliptic curve point",
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Deprecated:  "The groups_claim field is deprecated and will be removed in a future version. Use Authorization Server Claims (okta_auth_server_claim) or app profile configuration instead.",
				Description: "Groups claim for an OpenID Connect client application (DEPRECATED: This field will be removed in a future version. Use Authorization Server Claims instead).",
				Elem:        groupsClaimResource,
			},
			"app_settings_json": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Application settings in JSON format",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        utils.NormalizeDataJSON,
				DiffSuppressFunc: utils.NoChangeInObjectFromUnmarshaledJSON,
			},
			"authentication_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the associated app_signon_policy. If this property is removed from the application, the default sign-on-policy will be associated with this application. From now on, there is no need to attach authentication_policy for applications of type SERVICE`,
			},
			"jwks_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL reference to JWKS",
			},
			"skip_authentication_policy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Skip authentication policy operations. When set to true, the provider will not attempt to create, update, or delete authentication policies for this application.",
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

func resourceAppOAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV6ClientFromMetadata(meta)
	if err := validateGrantTypes(d); err != nil {
		return diag.Errorf("failed to create OAuth application: %v", err)
	}
	if err := validateAppOAuth(d, meta); err != nil {
		return diag.Errorf("failed to create OAuth application: %v", err)
	}
	app, err := buildAppOAuthV6(d, true)
	if err != nil {
		return diag.Errorf("failed to build OAuth application: %v", err)
	}
	activate := d.Get("status").(string) == StatusActive

	appResp, _, err := client.ApplicationAPI.CreateApplication(ctx).Application(app).Activate(activate).Execute()
	if err != nil {
		return diag.Errorf("failed to create OAuth application: %v", err)
	}

	oidcApp, err := verifyOidcAppTypeV6(*appResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(oidcApp.GetId())
	if !d.Get("omit_secret").(bool) && oidcApp.Credentials.OauthClient.HasClientSecret() {
		_ = d.Set("client_secret", oidcApp.Credentials.OauthClient.GetClientSecret())
	}

	err = handleAppLogo(ctx, d, meta, oidcApp.GetId(), oidcApp.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for OAuth application: %v", err)
	}

	err = setAppOauthGroupsClaim(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to update groups claim for an OAuth application: %v", err)
	}

	if d.Get("type") != "service" {
		if !d.Get("skip_authentication_policy").(bool) {
			err = createOrUpdateAuthenticationPolicy(ctx, d, meta, app.Id)
			if err != nil {
				return diag.Errorf("failed to set authentication policy for an OAuth application: %v", err)
			}
		}
	}

	return resourceAppOAuthRead(ctx, d, meta)
}

func setAppOauthGroupsClaim(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	raw, ok := d.GetOk("groups_claim")
	if !ok {
		return nil
	}

	// Log deprecation warning
	logger(meta).Warn("groups_claim is deprecated and will be removed in a future version. Please use Authorization Server Claims (okta_auth_server_claim) or app profile configuration instead.")

	c := meta.(*config.Config)
	if c.IsOAuth20Auth() {
		logger(meta).Warn("setting groups_claim disabled with OAuth 2.0 API authentication")
		return nil
	}

	// For now, keep the old behavior but with warnings
	// TODO: Remove in future version - functionality temporarily maintained for backward compatibility
	apiSupplement := getAPISupplementFromMetadata(meta)
	appID := d.Id()

	groupsClaim := raw.([]interface{})[0].(map[string]interface{})
	gc := buildGroupsClaimFromResource(groupsClaim)

	// Set issuer mode from the resource data for groups_claim API compatibility
	if issuerMode := d.Get("issuer_mode").(string); issuerMode != "" {
		gc.IssuerMode = issuerMode
	}

	_, err := apiSupplement.UpdateAppOauthGroupsClaim(ctx, appID, gc)
	return err
}

func updateAppOauthGroupsClaim(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("groups_claim") {
		return nil
	}

	// Log deprecation warning
	logger(meta).Warn("groups_claim is deprecated and will be removed in a future version. Please use Authorization Server Claims (okta_auth_server_claim) or app profile configuration instead.")

	return setAppOauthGroupsClaim(ctx, d, meta)
}

func buildGroupsClaimFromResource(groupsClaim map[string]interface{}) *sdk.AppOauthGroupClaim {
	gc := &sdk.AppOauthGroupClaim{
		Name:  groupsClaim["name"].(string),
		Value: groupsClaim["value"].(string),
	}

	// Map Terraform 'type' to API 'valueType'
	claimType := groupsClaim["type"].(string)
	switch claimType {
	case "FILTER":
		gc.ValueType = "GROUPS"
	case "EXPRESSION":
		gc.ValueType = "EXPRESSION"
	default:
		// Default to the provided type if it's a valid API value
		gc.ValueType = claimType
	}

	if filterType, ok := groupsClaim["filter_type"]; ok && filterType.(string) != "" {
		gc.GroupFilterType = filterType.(string)
	}

	return gc
}

func resourceAppOAuthRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV6ClientFromMetadata(meta)

	appResp, _, err := client.ApplicationAPI.GetApplication(ctx, d.Id()).Execute()
	if err != nil {
		return diag.Errorf("failed to get OAuth application: %v", err)
	}

	app, err := verifyOidcAppTypeV6(*appResp)
	if err != nil {
		return diag.Errorf("failed to verify app type: %v", err)
	}

	if app.GetId() == "" {
		d.SetId("")
		return nil
	}
	if !d.Get("skip_authentication_policy").(bool) {
		setAuthenticationPolicy(ctx, meta, d, app.Links)
	}
	var rawProfile string
	if profile := app.GetProfile(); profile != nil {
		p, _ := json.Marshal(profile)
		rawProfile = string(p)
	}

	credentials := app.GetCredentials()
	if userTemplate := credentials.GetUserNameTemplate(); userTemplate.GetTemplate() != "" {
		_ = d.Set("user_name_template", userTemplate.GetTemplate())
		_ = d.Set("user_name_template_type", userTemplate.GetType())
		_ = d.Set("user_name_template_suffix", userTemplate.GetUserSuffix())
		_ = d.Set("user_name_template_push_status", userTemplate.GetPushStatus())
	}

	settings := app.GetSettings()
	// Set basic app properties
	_ = d.Set("name", app.GetName())
	_ = d.Set("status", app.GetStatus())
	_ = d.Set("sign_on_mode", app.GetSignOnMode())
	_ = d.Set("label", app.GetLabel())

	// Set accessibility properties
	accessibility := app.GetAccessibility()
	_ = d.Set("accessibility_self_service", accessibility.GetSelfService())
	_ = d.Set("accessibility_error_redirect_url", accessibility.GetErrorRedirectUrl())
	_ = d.Set("accessibility_login_redirect_url", accessibility.GetLoginRedirectUrl())

	// Set visibility properties
	visibility := app.GetVisibility()
	_ = d.Set("auto_submit_toolbar", visibility.GetAutoSubmitToolbar())
	hide := visibility.GetHide()
	_ = d.Set("hide_ios", hide.GetIOS())
	_ = d.Set("hide_web", hide.GetWeb())

	// Set notes
	notes := settings.GetNotes()
	_ = d.Set("admin_note", notes.GetAdmin())
	_ = d.Set("enduser_note", notes.GetEnduser())
	_ = d.Set("profile", rawProfile)

	// Not setting client_secret, it is only provided on create and update for auth methods that require it
	oauthClient := credentials.GetOauthClient()
	if oauthClient.GetClientId() != "" {
		_ = d.Set("client_id", oauthClient.GetClientId())
		_ = d.Set("token_endpoint_auth_method", oauthClient.GetTokenEndpointAuthMethod())
		_ = d.Set("auto_key_rotation", oauthClient.GetAutoKeyRotation())
		_ = d.Set("pkce_required", oauthClient.GetPkceRequired())
	}

	// Handle app settings from AdditionalProperties
	if len(settings.AdditionalProperties) > 0 {
		appSettingsJSON, _ := json.Marshal(settings.AdditionalProperties)
		_ = d.Set("app_settings_json", string(appSettingsJSON))
	}

	_ = d.Set("logo_url", utils.LinksValue(app.Links, "logo", "href"))
	_ = d.Set("implicit_assignment", settings.GetImplicitAssignment())

	// Handle groups_claim with deprecation warning
	c := meta.(*config.Config)
	if c.IsOAuth20Auth() {
		logger(meta).Warn("reading groups_claim disabled with OAuth 2.0 API authentication")
	} else {
		gc, err := flattenGroupsClaim(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(gc) > 0 {
			logger(meta).Warn("groups_claim is deprecated and will be removed in a future version. Please use Authorization Server Claims (okta_auth_server_claim) or app profile configuration instead.")
		}
		_ = d.Set("groups_claim", gc)
	}

	return setOAuthClientSettingsV6(d, settings.OauthClient)
}

func flattenGroupsClaim(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]interface{}, error) {
	apiSupplement := getAPISupplementFromMetadata(meta)

	gc, _, err := apiSupplement.GetAppOauthGroupsClaim(ctx, d.Id())
	if err != nil {
		// If groups claim doesn't exist, return empty list rather than error
		return []interface{}{}, nil
	}

	if gc == nil || gc.Name == "" {
		return []interface{}{}, nil
	}

	// Map API 'valueType' back to Terraform 'type'
	terraformType := gc.ValueType
	if gc.ValueType == "GROUPS" && gc.GroupFilterType != "" {
		terraformType = "FILTER"
	}

	groupsClaimMap := map[string]interface{}{
		"type":        terraformType,
		"name":        gc.Name,
		"value":       gc.Value,
		"issuer_mode": gc.IssuerMode,
	}

	if gc.GroupFilterType != "" {
		groupsClaimMap["filter_type"] = gc.GroupFilterType
	}

	return []interface{}{groupsClaimMap}, nil
}

func setOAuthClientSettingsV6(d *schema.ResourceData, oauthClient *v6okta.OpenIdConnectApplicationSettingsClient) diag.Diagnostics {
	if oauthClient == nil {
		return nil
	}

	_ = d.Set("type", oauthClient.GetApplicationType())
	_ = d.Set("client_uri", oauthClient.GetClientUri())
	_ = d.Set("logo_uri", oauthClient.GetLogoUri())
	_ = d.Set("tos_uri", oauthClient.GetTosUri())
	_ = d.Set("policy_uri", oauthClient.GetPolicyUri())
	_ = d.Set("login_uri", oauthClient.GetInitiateLoginUri())
	_ = d.Set("jwks_uri", oauthClient.GetJwksUri())
	_ = d.Set("wildcard_redirect", oauthClient.GetWildcardRedirect())
	_ = d.Set("consent_method", oauthClient.GetConsentMethod())
	_ = d.Set("issuer_mode", oauthClient.GetIssuerMode())
	_ = d.Set("participate_slo", oauthClient.GetParticipateSlo())
	_ = d.Set("frontchannel_logout_uri", oauthClient.GetFrontchannelLogoutUri())
	_ = d.Set("frontchannel_logout_session_required", oauthClient.GetFrontchannelLogoutSessionRequired())

	if refreshToken := oauthClient.GetRefreshToken(); refreshToken.GetRotationType() != "" {
		_ = d.Set("refresh_token_rotation", refreshToken.GetRotationType())
		if refreshToken.HasLeeway() {
			_ = d.Set("refresh_token_leeway", int(refreshToken.GetLeeway()))
		}
	}

	// Handle response and grant types
	respTypes := make([]string, len(oauthClient.ResponseTypes))
	for i, responseType := range oauthClient.ResponseTypes {
		respTypes[i] = string(responseType)
	}

	grantTypes := make([]string, len(oauthClient.GrantTypes))
	for i, grantType := range oauthClient.GrantTypes {
		grantTypes[i] = string(grantType)
	}

	aggMap := map[string]interface{}{
		"redirect_uris":             oauthClient.RedirectUris,
		"response_types":            utils.ConvertStringSliceToSet(respTypes),
		"grant_types":               utils.ConvertStringSliceToSet(grantTypes),
		"post_logout_redirect_uris": utils.ConvertStringSliceToSet(oauthClient.PostLogoutRedirectUris),
	}

	if idpLogin := oauthClient.GetIdpInitiatedLogin(); idpLogin.GetMode() != "" {
		_ = d.Set("login_mode", idpLogin.GetMode())
		aggMap["login_scopes"] = utils.ConvertStringSliceToSet(idpLogin.DefaultScope)
	}

	err := utils.SetNonPrimitives(d, aggMap)
	if err != nil {
		return diag.Errorf("failed to set OAuth application properties: %v", err)
	}
	return nil
}

func resourceAppOAuthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	additionalChanges, err := AppUpdateStatus(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if !additionalChanges {
		return nil
	}

	client := getOktaV6ClientFromMetadata(meta)
	if err = validateGrantTypes(d); err != nil {
		return diag.Errorf("failed to update OAuth application: %v", err)
	}
	if err = validateAppOAuth(d, meta); err != nil {
		return diag.Errorf("failed to update OAuth application: %v", err)
	}

	app, err := buildAppOAuthV6(d, false)
	if err != nil {
		return diag.Errorf("failed to build OAuth application: %v", err)
	}

	// When omit_secret is true on update, we make sure that do not include
	// the client secret value in the api call.
	// This is to ensure that when this is "toggled on", the apply which this occurs also does
	// not do a final "reset" of the client secret value to the original stored in state.
	if d.Get("omit_secret").(bool) {
		// Get the underlying OpenIdConnectApplication and modify the credentials
		if oidcApp := app.OpenIdConnectApplication; oidcApp != nil {
			credentials := oidcApp.GetCredentials()
			oauthClient := credentials.GetOauthClient()
			oauthClient.SetClientSecret("")
			credentials.SetOauthClient(oauthClient)
			oidcApp.SetCredentials(credentials)
		}
	}

	appResp, _, err := client.ApplicationAPI.ReplaceApplication(ctx, d.Id()).Application(app).Execute()
	if err != nil {
		return diag.Errorf("failed to update OAuth application: %v", err)
	}

	updatedApp, err := verifyOidcAppTypeV6(*appResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The `client_secret` value is always returned from the API on update
	// regardless if we pass a value or not.
	// We need to make sure that we set the value in state based upon the `omit_secret` behavior
	// When `true`: We blank out the secret value
	// When `false`: We set the secret value to the value returned from the API
	if d.Get("omit_secret").(bool) {
		_ = d.Set("client_secret", "")
	} else if updatedApp.Credentials.OauthClient.HasClientSecret() {
		_ = d.Set("client_secret", updatedApp.Credentials.OauthClient.GetClientSecret())
	}

	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, meta, updatedApp.GetId(), updatedApp.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for OAuth application: %v", err)
		}
	}
	err = updateAppOauthGroupsClaim(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to update groups claim for an OAuth application: %v", err)
	}
	if !d.Get("skip_authentication_policy").(bool) {
		err = createOrUpdateAuthenticationPolicy(ctx, d, meta, app.Id)
		if err != nil {
			return diag.Errorf("failed to set authentication policy an OAuth application: %v", err)
		}
	}
	return resourceAppOAuthRead(ctx, d, meta)
}

func resourceAppOAuthDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to delete OAuth application: %v", err)
	}
	return nil
}

func buildAppOAuthV6(d *schema.ResourceData, isNew bool) (v6okta.ListApplications200ResponseInner, error) {
	app := v6okta.NewOpenIdConnectApplicationWithDefaults()
	appType := d.Get("type").(string)
	grantTypes := utils.ConvertInterfaceToStringSet(d.Get("grant_types"))
	responseTypes := utils.ConvertInterfaceToStringSetNullable(d.Get("response_types"))

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

		if utils.ContainsOne(grantTypes, implicit, clientCredentials) {
			responseTypes = append(responseTypes, "token")
		}

		if utils.ContainsOne(grantTypes, password, authorizationCode, refreshToken) {
			responseTypes = append(responseTypes, "code")
		}
	}

	app.SetLabel(d.Get("label").(string))
	app.SetName("oidc_client")
	app.SetSignOnMode("OPENID_CONNECT")

	// Build credentials
	authMethod := d.Get("token_endpoint_auth_method").(string)
	oauthClient := v6okta.NewApplicationCredentialsOAuthClientWithDefaults()
	if autoRotate, ok := d.GetOk("auto_key_rotation"); ok {
		oauthClient.SetAutoKeyRotation(autoRotate.(bool))
	}
	if clientID := d.Get("client_id").(string); clientID != "" {
		oauthClient.SetClientId(clientID)
	}
	oauthClient.SetTokenEndpointAuthMethod(authMethod)
	if clientSecret := d.Get("client_secret").(string); clientSecret != "" {
		oauthClient.SetClientSecret(clientSecret)
	}

	// Handle PKCE requirements
	var pkceRequired bool
	pkceVal := d.GetRawConfig().GetAttr("pkce_required")
	if pkceVal.IsNull() {
		if authMethod == "none" {
			pkceRequired = true
		} else if isNew && (appType == "native" || appType == "browser") {
			pkceRequired = true
		} else {
			pkceRequired = false
		}
	} else {
		switch {
		case pkceVal.True():
			pkceRequired = true
		case pkceVal.False():
			pkceRequired = false
		}
	}
	oauthClient.SetPkceRequired(pkceRequired)

	if sec, ok := d.GetOk("client_basic_secret"); ok {
		oauthClient.SetClientSecret(sec.(string))
	}

	credentials := v6okta.NewOAuthApplicationCredentialsWithDefaults()
	credentials.SetOauthClient(*oauthClient)
	credentials.SetUserNameTemplate(*BuildUserNameTemplateV6(d))
	app.SetCredentials(*credentials)

	// Build OAuth client settings
	oauthClientSettings := v6okta.NewOpenIdConnectApplicationSettingsClientWithDefaults()
	oauthClientSettings.SetApplicationType(appType)

	// Convert grant and response types to v6 format
	v6GrantTypes := make([]string, len(grantTypes))
	copy(v6GrantTypes, grantTypes)

	v6ResponseTypes := make([]string, len(responseTypes))
	copy(v6ResponseTypes, responseTypes)

	oauthClientSettings.SetGrantTypes(v6GrantTypes)
	oauthClientSettings.SetResponseTypes(v6ResponseTypes)

	if clientURI := d.Get("client_uri").(string); clientURI != "" {
		oauthClientSettings.SetClientUri(clientURI)
	}
	if consentMethod := d.Get("consent_method").(string); consentMethod != "" {
		oauthClientSettings.SetConsentMethod(consentMethod)
	}

	if redirectUris := utils.ConvertInterfaceToStringArr(d.Get("redirect_uris")); len(redirectUris) > 0 {
		oauthClientSettings.SetRedirectUris(redirectUris)
	}

	if postLogoutUris := utils.ConvertInterfaceToStringSetNullable(d.Get("post_logout_redirect_uris")); len(postLogoutUris) > 0 {
		oauthClientSettings.SetPostLogoutRedirectUris(postLogoutUris)
	}

	if loginURI := d.Get("login_uri").(string); loginURI != "" {
		oauthClientSettings.SetInitiateLoginUri(loginURI)
	}
	if logoURI := d.Get("logo_uri").(string); logoURI != "" {
		oauthClientSettings.SetLogoUri(logoURI)
	}
	if policyURI := d.Get("policy_uri").(string); policyURI != "" {
		oauthClientSettings.SetPolicyUri(policyURI)
	}
	if tosURI := d.Get("tos_uri").(string); tosURI != "" {
		oauthClientSettings.SetTosUri(tosURI)
	}
	if issuerMode := d.Get("issuer_mode").(string); issuerMode != "" {
		oauthClientSettings.SetIssuerMode(issuerMode)
	}
	if wildcardRedirect := d.Get("wildcard_redirect").(string); wildcardRedirect != "" {
		oauthClientSettings.SetWildcardRedirect(wildcardRedirect)
	}
	if jwksURI := d.Get("jwks_uri").(string); jwksURI != "" {
		oauthClientSettings.SetJwksUri(jwksURI)
	}
	if participateSlo, ok := d.GetOk("participate_slo"); ok {
		oauthClientSettings.SetParticipateSlo(participateSlo.(bool))
	}
	if frontchannelLogoutURI := d.Get("frontchannel_logout_uri").(string); frontchannelLogoutURI != "" {
		oauthClientSettings.SetFrontchannelLogoutUri(frontchannelLogoutURI)
	}
	if frontchannelLogoutSessionRequired, ok := d.GetOk("frontchannel_logout_session_required"); ok {
		oauthClientSettings.SetFrontchannelLogoutSessionRequired(frontchannelLogoutSessionRequired.(bool))
	}

	// Handle JWKS if provided
	if jwksList, ok := d.GetOk("jwks"); ok {
		jwksData := jwksList.([]interface{})
		if len(jwksData) > 0 {
			// For now, temporarily disable inline JWKS due to v6 SDK schema conflicts
			// This is a known limitation with the v6 SDK JWKS unmarshalling
			// Users should use jwks_uri instead for v6 SDK compatibility

			// TODO: Re-enable when v6 SDK JWKS schema conflict is resolved
			// For now, we'll skip setting JWKS to avoid the unmarshalling error
			/*
				jwks := v6okta.NewOpenIdConnectApplicationSettingsClientKeysWithDefaults()

				var keyData []interface{}
				for _, jwk := range jwksData {
					jwkMap := jwk.(map[string]interface{})
					key := map[string]interface{}{
						"kid": jwkMap["kid"].(string),
						"kty": jwkMap["kty"].(string),
					}
					if e, ok := jwkMap["e"]; ok && e.(string) != "" {
						key["e"] = e.(string)
					}
					if n, ok := jwkMap["n"]; ok && n.(string) != "" {
						key["n"] = n.(string)
					}
					if x, ok := jwkMap["x"]; ok && x.(string) != "" {
						key["x"] = x.(string)
					}
					if y, ok := jwkMap["y"]; ok && y.(string) != "" {
						key["y"] = y.(string)
					}
					keyData = append(keyData, key)
				}

				jwks.AdditionalProperties = map[string]interface{}{"keys": keyData}
				oauthClientSettings.SetJwks(*jwks)
			*/
		}
	}

	// Handle IDP initiated login
	if loginScopes := utils.ConvertInterfaceToStringSet(d.Get("login_scopes")); len(loginScopes) > 0 {
		idpLogin := v6okta.NewOpenIdConnectApplicationIdpInitiatedLoginWithDefaults()
		idpLogin.SetDefaultScope(loginScopes)
		if loginMode := d.Get("login_mode").(string); loginMode != "" {
			idpLogin.SetMode(loginMode)
		}
		oauthClientSettings.SetIdpInitiatedLogin(*idpLogin)
	}

	// Handle refresh token settings
	hasRefresh := false
	for _, grantType := range grantTypes {
		if grantType == refreshToken {
			hasRefresh = true
			break
		}
	}

	if hasRefresh {
		refresh := v6okta.NewOpenIdConnectApplicationSettingsRefreshTokenWithDefaults()
		if rotate, ok := d.GetOk("refresh_token_rotation"); ok {
			refresh.SetRotationType(rotate.(string))
		}
		if leeway, ok := d.GetOk("refresh_token_leeway"); ok {
			refresh.SetLeeway(int32(leeway.(int)))
		} else {
			refresh.SetLeeway(0)
		}
		oauthClientSettings.SetRefreshToken(*refresh)
	}

	// Build main settings
	settings := v6okta.NewOpenIdConnectApplicationSettingsWithDefaults()
	if implicitAssignment, ok := d.GetOk("implicit_assignment"); ok {
		settings.SetImplicitAssignment(implicitAssignment.(bool))
	}
	settings.SetOauthClient(*oauthClientSettings)
	settings.SetNotes(*BuildAppNotesV6(d))

	// Handle app settings
	if appSettingsJSON, ok := d.GetOk("app_settings_json"); ok {
		var appSettings map[string]interface{}
		if err := json.Unmarshal([]byte(appSettingsJSON.(string)), &appSettings); err != nil {
			return v6okta.ListApplications200ResponseInner{}, fmt.Errorf("failed to unmarshal app_settings_json: %w", err)
		}
		settings.AdditionalProperties = appSettings
	}

	app.SetSettings(*settings)
	visibility, err := BuildAppVisibilityV6(d)
	if err != nil {
		return v6okta.ListApplications200ResponseInner{}, err
	}
	app.SetVisibility(*visibility)
	app.SetAccessibility(*BuildAppAccessibilityV6(d))

	// Handle profile
	if rawAttrs, ok := d.GetOk("profile"); ok {
		var attrs map[string]interface{}
		str := rawAttrs.(string)
		if err := json.Unmarshal([]byte(str), &attrs); err != nil {
			return v6okta.ListApplications200ResponseInner{}, fmt.Errorf("failed to unmarshal profile: %w", err)
		}
		app.SetProfile(attrs)
	}

	return v6okta.OpenIdConnectApplicationAsListApplications200ResponseInner(app), nil
}

func validateGrantTypes(d *schema.ResourceData) error {
	grantTypeList := utils.ConvertInterfaceToStringSet(d.Get("grant_types"))
	// If grant_types are not set, we default to the bare minimum in func buildAppOAuth
	if len(grantTypeList) < 1 {
		return nil
	}
	appType := d.Get("type").(string)
	appMap, ok := appRequirementsByType[appType]
	if !ok {
		return nil
	}
	// There is some conditional validation around grant types depending on application type.
	return utils.ConditionalValidator("grant_types", appType, appMap.RequiredGrantTypes, appMap.ValidGrantTypes, grantTypeList)
}

func validateAppOAuth(d *schema.ResourceData, meta interface{}) error {
	// Handle groups_claim validation with deprecation warning
	raw, ok := d.GetOk("groups_claim")
	if ok {
		logger(meta).Warn("groups_claim is deprecated and will be removed in a future version. Please use Authorization Server Claims (okta_auth_server_claim) or app profile configuration instead.")

		c := meta.(*config.Config)
		if c.IsOAuth20Auth() {
			logger(meta).Warn("groups_claim arguments are disabled with OAuth 2.0 API authentication")
		} else {
			groupsClaim := raw.([]interface{})[0].(map[string]interface{})
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
		if d.Get("login_mode").(string) == "OKTA" && len(utils.ConvertInterfaceToStringSet(d.Get("login_scopes"))) < 1 {
			return errors.New("you have to set up non-empty 'login_scopes' when 'login_mode' is 'OKTA'")
		}
	}
	grantTypes := utils.ConvertInterfaceToStringSet(d.Get("grant_types"))
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
	for _, v := range utils.ConvertInterfaceToStringSetNullable(d.Get("response_types")) {
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

func verifyOidcAppTypeV6(app v6okta.ListApplications200ResponseInner) (*v6okta.OpenIdConnectApplication, error) {
	if app.OpenIdConnectApplication != nil {
		return app.OpenIdConnectApplication, nil
	}
	return nil, fmt.Errorf("unexpected application response type from Okta - not an OpenIdConnectApplication")
}

func BuildUserNameTemplateV6(d *schema.ResourceData) *v6okta.ApplicationCredentialsUsernameTemplate {
	template := v6okta.NewApplicationCredentialsUsernameTemplateWithDefaults()
	if tmpl := d.Get("user_name_template").(string); tmpl != "" {
		template.SetTemplate(tmpl)
	}
	if tmplType := d.Get("user_name_template_type").(string); tmplType != "" {
		template.SetType(tmplType)
	}
	if suffix := d.Get("user_name_template_suffix").(string); suffix != "" {
		template.SetUserSuffix(suffix)
	}
	if pushStatus := d.Get("user_name_template_push_status").(string); pushStatus != "" {
		template.SetPushStatus(pushStatus)
	}
	return template
}

func BuildAppAccessibilityV6(d *schema.ResourceData) *v6okta.ApplicationAccessibility {
	accessibility := v6okta.NewApplicationAccessibilityWithDefaults()
	if selfService, ok := d.GetOk("accessibility_self_service"); ok {
		accessibility.SetSelfService(selfService.(bool))
	}
	if errorURL := d.Get("accessibility_error_redirect_url").(string); errorURL != "" {
		accessibility.SetErrorRedirectUrl(errorURL)
	}
	if loginURL := d.Get("accessibility_login_redirect_url").(string); loginURL != "" {
		accessibility.SetLoginRedirectUrl(loginURL)
	}
	return accessibility
}

func BuildAppVisibilityV6(d *schema.ResourceData) (*v6okta.ApplicationVisibility, error) {
	visibility := v6okta.NewApplicationVisibilityWithDefaults()
	if autoSubmit, ok := d.GetOk("auto_submit_toolbar"); ok {
		visibility.SetAutoSubmitToolbar(autoSubmit.(bool))
	}

	hide := v6okta.NewApplicationVisibilityHideWithDefaults()
	if hideIOS, ok := d.GetOk("hide_ios"); ok {
		hide.SetIOS(hideIOS.(bool))
	}
	if hideWeb, ok := d.GetOk("hide_web"); ok {
		hide.SetWeb(hideWeb.(bool))
	}
	visibility.SetHide(*hide)

	if appLinks, ok := d.GetOk("app_links_json"); ok {
		var links map[string]bool
		if err := json.Unmarshal([]byte(appLinks.(string)), &links); err != nil {
			return nil, fmt.Errorf("failed to unmarshal app_links_json: %w", err)
		}
		visibility.SetAppLinks(links)
	}
	return visibility, nil
}

func BuildAppNotesV6(d *schema.ResourceData) *v6okta.ApplicationSettingsNotes {
	notes := v6okta.NewApplicationSettingsNotesWithDefaults()
	if admin, ok := d.GetOk("admin_note"); ok {
		notes.SetAdmin(admin.(string))
	}
	if enduser, ok := d.GetOk("enduser_note"); ok {
		notes.SetEnduser(enduser.(string))
	}
	return notes
}
