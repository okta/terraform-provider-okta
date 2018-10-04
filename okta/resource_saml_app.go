package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func resourceSamlApp() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			return nil
		},
		Create: resourceSamlAppCreate,
		Read:   resourceSamlAppRead,
		Update: resourceSamlAppUpdate,
		Delete: resourceSamlAppDelete,
		Exists: resourceSamlAppExists,
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Single Sign On URL",
				ValidateFunc: validateIsURL,
			},
			"sso_url_override": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Single Sign On URL override",
				ValidateFunc: validateIsURL,
			},
			"recipient": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The location where the app may present the SAML assertion",
				ValidateFunc: validateIsURL,
			},
			"recipient_override": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Recipient override URL",
				ValidateFunc: validateIsURL,
			},
			"destination": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Identifies the location where the SAML response is intended to be sent inside of the SAML assertion",
				ValidateFunc: validateIsURL,
			},
			"destination_override": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Destination URL override",
				ValidateFunc: validateIsURL,
			},
			"audience": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Audience URI",
			},
			"audience_override": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Audience URI override",
			},
			"idp_issuer": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SAML issuer ID",
			},
			"subject_name_id_template": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template for app user's username when a user is assigned to the app",
			},
			"subject_name_id_format": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies the SAML processing rules.",
			},
			"response_signed": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether the SAML auth response message is digitally signed",
			},
			"assertion_signed": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether the SAML assertion is digitally signed",
			},
			"signature_algorithm": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Signature algorithm used ot digitally sign the assertion and response",
				ValidateFunc: validation.StringInSlice([]string{"RSA-SHA256", "RSA-SHA1"}, false),
			},
			"digest_algorithm": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Determines the digest algorithm used to digitally sign the SAML assertion and response",
				ValidateFunc: validation.StringInSlice([]string{"SHA256", "SHA1"}, false),
			},
			"honor_force_authn": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Prompt user to re-authenticate if SP asks for it",
			},
			"authn_context_class_ref": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies the SAML authentication context class for the assertion’s authentication statement",
			},
			"attribute_statements": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "SAML attributes. https://developer.okta.com/docs/api/resources/apps#attribute-statements-object",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The reference name of the attribute statement.",
						},
						"namespace": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name format of the attribute.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "The value of the attribute.",
							Elem: &schema.Schema{
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
			},
		},
	}
}

func resourceSamlAppExists(d *schema.ResourceData, m interface{}) (bool, error) {
	app := okta.NewSamlApplication()
	err := fetchApp(d, m, app)

	// Not sure if a non-nil app with an empty ID is possible but checking to avoid false positives.
	return app != nil && app.Id != "", err
}

func resourceSamlAppCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildSamlApp(d, m)
	activate := d.Get("status").(string) == "ACTIVE"
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(app, params)

	if err != nil {
		return err
	}

	d.SetId(app.Id)

	return resourceSamlAppRead(d, m)
}

func resourceSamlAppRead(d *schema.ResourceData, m interface{}) error {
	app := okta.NewSamlApplication()
	err := fetchApp(d, m, app)

	if err != nil {
		return err
	}

	d.Set("name", app.Name)
	d.Set("status", app.Status)
	d.Set("sign_on_mode", app.SignOnMode)
	d.Set("label", app.Label)
	d.Set("default_relay_state", app.Settings.SignOn.DefaultRelayState)
	d.Set("sso_url", app.Settings.SignOn.SsoAcsUrl)
	d.Set("sso_url_override", app.Settings.SignOn.SsoAcsUrlOverride)
	d.Set("recipient", app.Settings.SignOn.Recipient)
	d.Set("recipient_override", app.Settings.SignOn.RecipientOverride)
	d.Set("destination", app.Settings.SignOn.Destination)
	d.Set("destination_override", app.Settings.SignOn.DestinationOverride)
	d.Set("audience", app.Settings.SignOn.Audience)
	d.Set("audience_override", app.Settings.SignOn.AudienceOverride)
	d.Set("idp_issuer", app.Settings.SignOn.IdpIssuer)
	d.Set("subject_name_id_template", app.Settings.SignOn.SubjectNameIdTemplate)
	d.Set("subject_name_id_format", app.Settings.SignOn.SubjectNameIdFormat)
	d.Set("response_signed", app.Settings.SignOn.ResponseSigned)
	d.Set("assertion_signed", app.Settings.SignOn.AssertionSigned)
	d.Set("signature_algorithm", app.Settings.SignOn.SignatureAlgorithm)
	d.Set("digest_algorithm", app.Settings.SignOn.DigestAlgorithm)
	d.Set("honor_force_authn", app.Settings.SignOn.HonorForceAuthn)
	d.Set("authn_context_class_ref", app.Settings.SignOn.AuthnContextClassRef)

	return nil
}

func resourceSamlAppUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildSamlApp(d, m)
	_, _, err := client.Application.UpdateApplication(d.Id(), app)

	if err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	err = setAppStatus(d, client, app.Status, desiredStatus)

	if err != nil {
		return err
	}

	return resourceSamlAppRead(d, m)
}

func resourceSamlAppDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.Application.DeactivateApplication(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(d.Id())

	return err
}

func buildSamlApp(d *schema.ResourceData, m interface{}) *okta.SamlApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewSamlApplication()

	responseSigned := d.Get("response_signed").(bool)
	assertionSigned := d.Get("assertion_signed").(bool)
	honorForce := d.Get("honor_force_authn").(bool)
	app.Label = d.Get("label").(string)
	app.Name = d.Get("name").(string)
	app.Settings = okta.NewSamlApplicationSettings()
	app.Settings.SignOn = &okta.SamlApplicationSettingsSignOn{
		DefaultRelayState:     d.Get("default_relay_state").(string),
		SsoAcsUrl:             d.Get("sso_url").(string),
		SsoAcsUrlOverride:     d.Get("sso_url_override").(string),
		Recipient:             d.Get("recipient").(string),
		Destination:           d.Get("destination").(string),
		DestinationOverride:   d.Get("destination_override").(string),
		Audience:              d.Get("audience").(string),
		AudienceOverride:      d.Get("audience_override").(string),
		IdpIssuer:             d.Get("idp_issuer").(string),
		SubjectNameIdTemplate: d.Get("subject_name_id_template").(string),
		SubjectNameIdFormat:   d.Get("subject_name_id_format").(string),
		ResponseSigned:        &responseSigned,
		AssertionSigned:       &assertionSigned,
		SignatureAlgorithm:    d.Get("signature_algorithm").(string),
		DigestAlgorithm:       d.Get("digest_algorithm").(string),
		HonorForceAuthn:       &honorForce,
		AuthnContextClassRef:  d.Get("authn_context_class_ref").(string),
	}
	app.Credentials = okta.NewApplicationCredentials()
	app.Credentials.Signing = okta.NewApplicationCredentialsSigning()
	app.Credentials.UserNameTemplate = okta.NewApplicationCredentialsUsernameTemplate()

	return app
}
