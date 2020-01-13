package okta

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
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
	"web": &applicationMap{
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
	"native": &applicationMap{
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
	"browser": &applicationMap{
		ValidGrantTypes: []string{
			implicit,
			authorizationCode,
		},
	},
	"service": &applicationMap{
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
		Create: resourceAppOAuthCreate,
		Read:   resourceAppOAuthRead,
		Update: resourceAppOAuthUpdate,
		Delete: resourceAppOAuthDelete,
		Exists: resourceAppOAuthExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			// Force new if omit_secret goes from true to false
			if d.Id() != "" {
				old, new := d.GetChange("omit_secret")
				if old.(bool) == true && new.(bool) == false {
					return d.ForceNew("omit_secret")
				}
			}
			return nil
		},
		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Schema: buildAppSchema(map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"web", "native", "browser", "service"}, false),
				Required:     true,
				ForceNew:     true,
				Description:  "The type of client application.",
			},
			"client_id": &schema.Schema{
				Type: schema.TypeString,
				// This field is Optional + Computed because we automatically set the
				// client_id value if none is specified during creation.
				// If the client_id is set after creation, the resource will be recreated only if its different from
				// the computed client_id.
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "OAuth client ID. If set during creation, app is created with this id.",
			},
			"custom_client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "This property allows you to set your client_id.",
				Deprecated:  "This field is being replaced by client_id. Please set that field instead.",
			},
			"omit_secret": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				// No ForceNew to avoid recreating when going from false => true
				Description: "This tells the provider not to persist the application's secret to state. If this is ever changes from true => false your app will be recreated.",
				Default:     false,
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "OAuth client secret key. This will be in plain text in your statefile unless you set omit_secret above.",
			},
			"client_basic_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "OAuth client secret key, this can be set when token_endpoint_auth_method is client_secret_basic.",
			},
			"token_endpoint_auth_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"none", "client_secret_post", "client_secret_basic", "client_secret_jwt"},
					false,
				),
				Default:     "client_secret_basic",
				Description: "Requested authentication method for the token endpoint.",
			},
			"auto_key_rotation": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Requested key rotation mode.",
			},
			"client_uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URI to a web page providing information about the client.",
			},
			"logo_uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URI that references a logo for the client.",
			},
			"login_uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URI that initiates login.",
			},
			"redirect_uris": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of URIs for use in the redirect-based flow. This is required for all application types except service. Note: see okta_app_oauth_redirect_uri for appending to this list in a decentralized way.",
			},
			"post_logout_redirect_uris": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of URIs for redirection after logout",
			},
			"response_types": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"code", "token", "id_token"}, false),
				},
				Optional:    true,
				Description: "List of OAuth 2.0 response type strings.",
			},
			"grant_types": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice(
						[]string{authorizationCode, implicit, password, refreshToken, clientCredentials},
						false,
					),
				},
				Optional:    true,
				Description: "List of OAuth 2.0 grant types. Conditional validation params found here https://developer.okta.com/docs/api/resources/apps#credentials-settings-details. Defaults to minimum requirements per app type.",
			},
			// "Early access" properties.. looks to be in beta which requires opt-in per account
			"tos_uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "*Early Access Property*. URI to web page providing client tos (terms of service).",
			},
			"policy_uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "*Early Access Property*. URI to web page providing client policy document.",
			},
			"consent_method": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"REQUIRED", "TRUSTED"}, false),
				Description:  "*Early Access Property*. Indicates whether user consent is required or implicit. Valid values: REQUIRED, TRUSTED. Default value is TRUSTED",
			},
			"issuer_mode": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CUSTOM_URL", "ORG_URL"}, false),
				Default:      "ORG_URL",
				Description:  "*Early Access Property*. Indicates whether the Okta Authorization Server uses the original Okta org domain URL or a custom domain URL as the issuer of ID token for this client.",
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
				Default:     true,
				Description: "Do not display application icon on mobile app",
			},
			"hide_web": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Do not display application icon to users",
			},
			"profile": &schema.Schema{
				Type:        schema.TypeString,
				StateFunc:   normalizeDataJSON,
				Optional:    true,
				Description: "Custom JSON that represents an OAuth application's profile",
			},
		}),
	}
}

func resourceAppOAuthExists(d *schema.ResourceData, m interface{}) (bool, error) {
	app := okta.NewOpenIdConnectApplication()
	err := fetchApp(d, m, app)

	// Not sure if a non-nil app with an empty ID is possible but checking to avoid false positives.
	return app != nil && app.Id != "", err
}

func validateGrantTypes(d *schema.ResourceData) error {
	grantTypeList := convertInterfaceToStringSet(d.Get("grant_types"))
	appType := d.Get("type").(string)
	appMap := appGrantTypeMap[appType]

	// There is some conditional validation around grant types depending on application type.
	return conditionalValidator("grant_types", appType, appMap.RequiredGrantTypes, appMap.ValidGrantTypes, grantTypeList)
}

// validateClientID returns an error if both custom_client_id and client_id are specified
// TODO: remove this in the next release when custom_client_id is marked as removed
func validateClientID(d *schema.ResourceData) error {
	if _, ok := d.GetOk("custom_client_id"); ok {
		if _, ok := d.GetOk("client_id"); ok {
			return fmt.Errorf("cannot specify both custom_client_id and client_id. Specify only client_id instead")
		}
	}
	return nil
}

func resourceAppOAuthCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	if err := validateGrantTypes(d); err != nil {
		return err
	}
	if err := validateClientID(d); err != nil {
		return err
	}

	app := buildAppOAuth(d, m)
	desiredStatus := d.Get("status").(string)
	activate := desiredStatus == "ACTIVE"
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(app, params)

	if err != nil {
		return err
	}

	d.SetId(app.Id)
	if !d.Get("omit_secret").(bool) {
		// Needs to be set immediately, not provided again after this
		d.Set("client_secret", app.Credentials.OauthClient.ClientSecret)
	}
	err = handleAppGroupsAndUsers(app.Id, d, m)

	if err != nil {
		return err
	}

	return resourceAppOAuthRead(d, m)
}

func resourceAppOAuthRead(d *schema.ResourceData, m interface{}) error {
	app := okta.NewOpenIdConnectApplication()
	err := fetchApp(d, m, app)

	if app == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("name", app.Name)
	d.Set("status", app.Status)
	d.Set("sign_on_mode", app.SignOnMode)
	d.Set("label", app.Label)
	d.Set("profile", app.Profile)
	d.Set("type", app.Settings.OauthClient.ApplicationType)
	// Not setting client_secret, it is only provided on create for auth methods that require it
	d.Set("client_id", app.Credentials.OauthClient.ClientId)
	d.Set("token_endpoint_auth_method", app.Credentials.OauthClient.TokenEndpointAuthMethod)
	d.Set("auto_key_rotation", app.Credentials.OauthClient.AutoKeyRotation)
	d.Set("consent_method", app.Settings.OauthClient.ConsentMethod)
	d.Set("client_uri", app.Settings.OauthClient.ClientUri)
	d.Set("logo_uri", app.Settings.OauthClient.LogoUri)
	d.Set("tos_uri", app.Settings.OauthClient.TosUri)
	d.Set("policy_uri", app.Settings.OauthClient.PolicyUri)
	d.Set("login_uri", app.Settings.OauthClient.InitiateLoginUri)
	d.Set("auto_submit_toolbar", app.Visibility.AutoSubmitToolbar)
	d.Set("hide_ios", app.Visibility.Hide.IOS)
	d.Set("hide_web", app.Visibility.Hide.Web)

	if app.Settings.OauthClient.IssuerMode != "" {
		d.Set("issuer_mode", app.Settings.OauthClient.IssuerMode)
	}

	// If this is ever changed omit it.
	if d.Get("omit_secret").(bool) {
		d.Set("client_secret", "")
	}

	if err = syncGroupsAndUsers(app.Id, d, m); err != nil {
		return err
	}
	aggMap := map[string]interface{}{
		"redirect_uris":             convertStringSetToInterface(app.Settings.OauthClient.RedirectUris),
		"response_types":            convertStringSetToInterface(app.Settings.OauthClient.ResponseTypes),
		"grant_types":               convertStringSetToInterface(app.Settings.OauthClient.GrantTypes),
		"post_logout_redirect_uris": convertStringSetToInterface(app.Settings.OauthClient.PostLogoutRedirectUris),
	}

	return setNonPrimitives(d, aggMap)
}

func resourceAppOAuthUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	if err := validateGrantTypes(d); err != nil {
		return err
	}
	if err := validateClientID(d); err != nil {
		return err
	}

	app := buildAppOAuth(d, m)
	if _, _, err := client.Application.UpdateApplication(d.Id(), app); err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	if err := setAppStatus(d, client, app.Status, desiredStatus); err != nil {
		return err
	}

	if err := handleAppGroupsAndUsers(app.Id, d, m); err != nil {
		return err
	}

	return resourceAppOAuthRead(d, m)
}

func resourceAppOAuthDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)

	if d.Get("status").(string) == "ACTIVE" {
		_, err := client.Application.DeactivateApplication(d.Id())
		if err != nil {
			return err
		}
	}

	_, err := client.Application.DeleteApplication(d.Id())
	return err
}

func buildAppOAuth(d *schema.ResourceData, m interface{}) *okta.OpenIdConnectApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewOpenIdConnectApplication()

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

	if ccid, ok := d.GetOk("custom_client_id"); ok {
		// if client_id is set, gets precedence, since client_id is the future source of truth of the applications
		// client_id.
		// if client_id is not set, set it to the value specified by custom_client_id.
		if cid, ok := d.GetOk("client_id"); !ok {
			d.Set("client_id", cid.(string))
			app.Credentials.OauthClient.ClientId = ccid.(string)
		}
	}

	app.Settings = &okta.OpenIdConnectApplicationSettings{
		OauthClient: &okta.OpenIdConnectApplicationSettingsClient{
			ApplicationType:        appType,
			ClientUri:              d.Get("client_uri").(string),
			ConsentMethod:          d.Get("consent_method").(string),
			GrantTypes:             grantTypes,
			InitiateLoginUri:       d.Get("login_uri").(string),
			LogoUri:                d.Get("logo_uri").(string),
			PolicyUri:              d.Get("policy_uri").(string),
			RedirectUris:           convertInterfaceToStringSetNullable(d.Get("redirect_uris")),
			PostLogoutRedirectUris: convertInterfaceToStringSetNullable(d.Get("post_logout_redirect_uris")),
			ResponseTypes:          responseTypes,
			TosUri:                 d.Get("tos_uri").(string),
			IssuerMode:             d.Get("issuer_mode").(string),
		},
	}
	app.Visibility = buildVisibility(d)

	if rawAttrs, ok := d.GetOk("profile"); ok {
		var attrs map[string]interface{}
		str := rawAttrs.(string)
		json.Unmarshal([]byte(str), &attrs)

		app.Profile = attrs
	}

	return app
}

func getClusterID() {}
