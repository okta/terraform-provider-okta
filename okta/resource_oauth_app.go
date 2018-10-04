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

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of resource.",
			},
			"sign_on_mode": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sign on mode of application.",
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Application label.",
			},
			"status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ACTIVE",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
				Description:  "Status of application.",
			},
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
			"redirect_uris": &schema.Schema{
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of URIs for use in the redirect-based flow. This is required for all application types except service.",
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
				Description: "*Beta Okta Property*. URI to web page providing client tos (terms of service).",
			},
			"policy_uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "*Beta Okta Property*. URI to web page providing client policy document.",
			},
			"consent_method": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"REQUIRED", "TRUSTED"}, false),
				Description:  "*Beta Okta Property*. Indicates whether user consent is required or implicit. Valid values: REQUIRED, TRUSTED. Default value is TRUSTED",
			},
		},
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

	return setNonPrimitives(d, map[string]interface{}{
		"redirect_uris":  app.Settings.OauthClient.RedirectUris,
		"response_types": app.Settings.OauthClient.ResponseTypes,
		"grant_types":    app.Settings.OauthClient.GrantTypes,
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

	return resourceOAuthAppRead(d, m)
}

func resourceOAuthAppDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.Application.DeactivateApplication(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(d.Id())

	return err
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
			ApplicationType: appType,
			ClientUri:       d.Get("client_uri").(string),
			ConsentMethod:   d.Get("consent_method").(string),
			GrantTypes:      grantTypes,
			LogoUri:         d.Get("logo_uri").(string),
			PolicyUri:       d.Get("policy_uri").(string),
			RedirectUris:    convertInterfaceToStringArrNullable(d.Get("redirect_uris")),
			ResponseTypes:   responseTypes,
			TosUri:          d.Get("tos_uri").(string),
		},
	}

	return app
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
