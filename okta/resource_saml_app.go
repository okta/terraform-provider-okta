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
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			grantTypeList := convertInterfaceToStringArr(d.Get("grant_types"))
			appType := d.Get("type").(string)
			appMap := appGrantTypeMap[appType]

			// There is some conditional validation around grant types depending on application type.
			return conditionalValidator("grant_types", appType, appMap.RequiredGrantTypes, appMap.ValidGrantTypes, grantTypeList)
		},
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
				Optional:    true,
				Description: "Name of preexisting SAML application.",
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
			"default_relay_state": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies a specific application resource in an IDP initiated SSO scenario.",
			},
			"sso_url": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Description: "Single Sign On URL",
				ValidateFunc: validateIsURL
			}, 
			"sso_url_override": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Description: "Single Sign On URL override",
				ValidateFunc: validateIsURL
			}, 
			"recipient": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Description: "The location where the app may present the SAML assertion",
				ValidateFunc: validateIsURL
			}, 
			"recipient_override": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Description: "Recipient override URL",
				ValidateFunc: validateIsURL
			}, 
			"destination": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Description: "Identifies the location where the SAML response is intended to be sent inside of the SAML assertion",
				ValidateFunc: validateIsURL
			}, 
		},
	}
}

func resourceOAuthAppExists(d *schema.ResourceData, m interface{}) (bool, error) {
	app, err := fetchApp(d, m)

	// Not sure if a non-nil app with an empty ID is possible but checking to avoid false positives.
	return app != nil && app.Id != "", err
}

func resourceOAuthAppCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
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
	app, err := fetchApp(d, m)

	if err != nil {
		return err
	}

	d.Set("name", app.Name)
	d.Set("status", app.Status)
	d.Set("sign_on_mode", app.SignOnMode)
	d.Set("label", app.Label)
}

func resourceOAuthAppUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildOAuthApp(d, m)
	_, _, err := client.Application.UpdateApplication(d.Id(), app)

	desiredStatus := d.Get("status").(string)

	if app.Status != desiredStatus {
		if desiredStatus == "INACTIVE" {
			client.Application.DeactivateApplication(d.Id())
		} else if desiredStatus == "ACTIVE" {
			client.Application.ActivateApplication(d.Id())
		}
	}

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

	return err
}

func fetchApp(d *schema.ResourceData, m interface{}) (*okta.SamlApplication, error) {
	client := getOktaClientFromMetadata(m)
	params := &query.Params{}
	newApp := &okta.SampleApplication{}
	_, response, err := client.Application.GetApplication(d.Id(), newApp, params)

	// We don't want to consider a 404 an error in some cases and thus the delineation
	if response.StatusCode == 404 {
		return nil, nil
	}

	return newApp, err
}

func buildOAuthApp(d *schema.ResourceData, m interface{}) *okta.SamlApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewSamlApplication()

	app.Label = d.Get("label").(string)
	app.Credentials = okta.ApplicationCredentials{}

	return app
}
