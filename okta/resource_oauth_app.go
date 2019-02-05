package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
		ValidGrantTypes: []string{implicit},
	},
	"service": &applicationMap{
		ValidGrantTypes: []string{clientCredentials},
	},
}

func resourceOAuthApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAuthAppCreate,
		Read:   resourceOAuthAppRead,
		Update: resourceOAuthAppUpdate,
		Delete: resourceOAuthAppDelete,
		Exists: resourceOAuthAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: buildAppSchema(map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"web", "native", "browser", "service"}, false),
				Required:     true,
				ForceNew:     true,
				Description:  "The type of client application.",
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OAuth client ID.",
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OAuth client secret key.",
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
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of URIs for use in the redirect-based flow. This is required for all application types except service.",
			},
			"post_logout_redirect_uris": &schema.Schema{
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of URIs for redirection after logout",
			},
			"response_types": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateFunc:     validation.StringInSlice([]string{"code", "token", "id_token"}, false),
					DiffSuppressFunc: suppressDefaultedDiff,
				},
				Optional:         true,
				DiffSuppressFunc: suppressDefaultedArrayDiff,
				Description:      "List of OAuth 2.0 response type strings.",
			},
			"grant_types": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice(
						[]string{authorizationCode, implicit, password, refreshToken, clientCredentials},
						false,
					),
					DiffSuppressFunc: suppressDefaultedDiff,
				},
				Optional:         true,
				Description:      "List of OAuth 2.0 grant types. Conditional validation params found here https://developer.okta.com/docs/api/resources/apps#credentials-settings-details. Defaults to minimum requirements per app type.",
				DiffSuppressFunc: suppressDefaultedArrayDiff,
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
				Description:  "*Early Access Property*. Indicates whether the Okta Authorization Server uses the original Okta org domain URL or a custom domain URL as the issuer of ID token for this client.",
			},
		}),
	}
}

func resourceOAuthAppExists(d *schema.ResourceData, m interface{}) (bool, error) {
	app := okta.NewOpenIdConnectApplication()
	err := fetchApp(d, m, app)

	// Not sure if a non-nil app with an empty ID is possible but checking to avoid false positives.
	return app != nil && app.Id != "", err
}

func validateGrantTypes(d *schema.ResourceData) error {
	grantTypeList := convertInterfaceToStringArr(d.Get("grant_types"))
	appType := d.Get("type").(string)
	appMap := appGrantTypeMap[appType]

	// There is some conditional validation around grant types depending on application type.
	return conditionalValidator("grant_types", appType, appMap.RequiredGrantTypes, appMap.ValidGrantTypes, grantTypeList)
}

func resourceOAuthAppCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	if err := validateGrantTypes(d); err != nil {
		return err
	}

	app := buildOAuthApp(d, m)
	desiredStatus := d.Get("status").(string)
	activate := desiredStatus == "ACTIVE"
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(app, params)

	if err != nil {
		return err
	}

	d.SetId(app.Id)
	err = handleAppGroupsAndUsers(app.Id, d, m)

	if err != nil {
		return err
	}

	return resourceOAuthAppRead(d, m)
}

func resourceOAuthAppRead(d *schema.ResourceData, m interface{}) error {
	app := okta.NewOpenIdConnectApplication()
	err := fetchApp(d, m, app)

	if err != nil {
		return err
	}

	d.Set("name", app.Name)
	d.Set("status", app.Status)
	d.Set("sign_on_mode", app.SignOnMode)
	d.Set("label", app.Label)
	d.Set("type", app.Settings.OauthClient.ApplicationType)
	d.Set("client_id", app.Credentials.OauthClient.ClientId)
	d.Set("client_secret", app.Credentials.OauthClient.ClientSecret)
	d.Set("token_endpoint_auth_method", app.Credentials.OauthClient.TokenEndpointAuthMethod)
	d.Set("auto_key_rotation", app.Credentials.OauthClient.AutoKeyRotation)
	d.Set("consent_method", app.Settings.OauthClient.ConsentMethod)
	d.Set("client_uri", app.Settings.OauthClient.ClientUri)
	d.Set("logo_uri", app.Settings.OauthClient.LogoUri)
	d.Set("tos_uri", app.Settings.OauthClient.TosUri)
	d.Set("policy_uri", app.Settings.OauthClient.PolicyUri)
	d.Set("login_uri", app.Settings.OauthClient.InitiateLoginUri)
	d.Set("issuer_mode", app.Settings.OauthClient.IssuerMode)

	if err = syncGroupsAndUsers(app.Id, d, m); err != nil {
		return err
	}

	return setNonPrimitives(d, map[string]interface{}{
		"redirect_uris":             app.Settings.OauthClient.RedirectUris,
		"response_types":            app.Settings.OauthClient.ResponseTypes,
		"grant_types":               app.Settings.OauthClient.GrantTypes,
		"post_logout_redirect_uris": app.Settings.OauthClient.PostLogoutRedirectUris,
	})
}

func resourceOAuthAppUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	if err := validateGrantTypes(d); err != nil {
		return err
	}

	app := buildOAuthApp(d, m)
	_, _, err := client.Application.UpdateApplication(d.Id(), app)
	if err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	err = setAppStatus(d, client, app.Status, desiredStatus)
	err = handleAppGroupsAndUsers(app.Id, d, m)

	if err != nil {
		return err
	}

	return resourceOAuthAppRead(d, m)
}

func resourceOAuthAppDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.Application.DeactivateApplication(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(d.Id())
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}

func buildOAuthApp(d *schema.ResourceData, m interface{}) *okta.OpenIdConnectApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewOpenIdConnectApplication()

	// Need to a bool pointer, it appears the Okta SDK uses this as a way to avoid false being omitted.
	keyRotation := d.Get("auto_key_rotation").(bool)
	appType := d.Get("type").(string)
	grantTypes := convertInterfaceToStringArr(d.Get("grant_types"))
	responseTypes := convertInterfaceToStringArrNullable(d.Get("response_types"))

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
	app.Credentials = &okta.OAuthApplicationCredentials{
		OauthClient: &okta.ApplicationCredentialsOAuthClient{
			AutoKeyRotation:         &keyRotation,
			ClientId:                d.Get("client_id").(string),
			ClientSecret:            d.Get("client_secret").(string),
			TokenEndpointAuthMethod: d.Get("token_endpoint_auth_method").(string),
		},
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
			RedirectUris:           convertInterfaceToStringArrNullable(d.Get("redirect_uris")),
			PostLogoutRedirectUris: convertInterfaceToStringArrNullable(d.Get("post_logout_redirect_uris")),
			ResponseTypes:          responseTypes,
			TosUri:                 d.Get("tos_uri").(string),
			IssuerMode:             d.Get("issuer_mode").(string),
		},
	}

	return app
}
