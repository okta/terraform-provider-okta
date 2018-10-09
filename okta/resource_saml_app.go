package okta

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

// Fields required if preconfigured_app is not provided
var customSamlAppRequiredFields = []string{
	"sso_url",
	"recipient",
	"destination",
	"audience",
	"idp_issuer",
	"subject_name_id_template",
	"subject_name_id_format",
	"signature_algorithm",
	"digest_algorithm",
	"honor_force_authn",
	"authn_context_class_ref",
}

func resourceSamlApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceSamlAppCreate,
		Read:   resourceSamlAppRead,
		Update: resourceSamlAppUpdate,
		Delete: resourceSamlAppDelete,
		Exists: resourceSamlAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"preconfigured_app": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of preexisting SAML application.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "App name.",
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
			"auto_submit_toolbar": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Display auto submit toolbar",
			},
			"hide_ios": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not display application icon on mobile app",
			},
			"hide_web": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not display application icon to users",
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
			"recipient": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The location where the app may present the SAML assertion",
				ValidateFunc: validateIsURL,
			},
			"destination": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Identifies the location where the SAML response is intended to be sent inside of the SAML assertion",
				ValidateFunc: validateIsURL,
			},
			"audience": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Audience URI",
				ValidateFunc: validateIsURL,
			},
			"idp_issuer": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SAML issuer ID",
			},
			"sp_issuer": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SAML SP issuer ID",
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
				ValidateFunc: validation.StringInSlice(
					[]string{
						"urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
						"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
						"urn:oasis:names:tc:SAML:1.1:nameid-format:x509SubjectName",
						"urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
						"urn:oasis:names:tc:SAML:2.0:nameid-format:transient",
					},
					false,
				),
			},
			"response_signed": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether the SAML auth response message is digitally signed",
			},
			"request_compressed": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Denotes whether the request is compressed or not.",
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
				ValidateFunc: validation.StringInSlice([]string{"RSA_SHA256", "RSA_SHA1"}, false),
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
			"accessibility_self_service": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable self service",
			},
			"accessibility_error_redirect_url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Custom error page URL",
				ValidateFunc: validateIsURL,
			},
			"accessibility_login_redirect_url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Custom login page URL",
				ValidateFunc: validateIsURL,
			},
		},
	}
}

func resourceSamlAppExists(d *schema.ResourceData, m interface{}) (bool, error) {
	app := okta.NewApplication()
	err := fetchApp(d, m, app)

	// Not sure if a non-nil app with an empty ID is possible but checking to avoid false positives.
	return app != nil && app.Id != "", err
}

func resourceSamlAppCreate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app, err := buildApp(d, m)

	if err != nil {
		return err
	}

	activate := d.Get("status").(string) == "ACTIVE"
	params := &query.Params{Activate: &activate}
	_, _, err = client.Application.CreateApplication(app, params)

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
	d.Set("recipient", app.Settings.SignOn.Recipient)
	d.Set("destination", app.Settings.SignOn.Destination)
	d.Set("audience", app.Settings.SignOn.Audience)
	d.Set("idp_issuer", app.Settings.SignOn.IdpIssuer)
	d.Set("subject_name_id_template", app.Settings.SignOn.SubjectNameIdTemplate)
	d.Set("subject_name_id_format", app.Settings.SignOn.SubjectNameIdFormat)
	d.Set("response_signed", app.Settings.SignOn.ResponseSigned)
	d.Set("assertion_signed", app.Settings.SignOn.AssertionSigned)
	d.Set("signature_algorithm", app.Settings.SignOn.SignatureAlgorithm)
	d.Set("digest_algorithm", app.Settings.SignOn.DigestAlgorithm)
	d.Set("honor_force_authn", app.Settings.SignOn.HonorForceAuthn)
	d.Set("authn_context_class_ref", app.Settings.SignOn.AuthnContextClassRef)
	d.Set("accessibility_self_service", app.Accessibility.SelfService)
	d.Set("accessibility_login_redirect_url", app.Accessibility.LoginRedirectUrl)
	d.Set("accessibility_error_redirect_url", app.Accessibility.ErrorRedirectUrl)

	return nil
}

func resourceSamlAppUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app, err := buildApp(d, m)

	if err != nil {
		return err
	}

	_, _, err = client.Application.UpdateApplication(d.Id(), app)

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

func buildApp(d *schema.ResourceData, m interface{}) (*okta.SamlApplication, error) {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewSamlApplication()
	app.Label = d.Get("label").(string)
	app.SignOnMode = "SAML_2_0"
	responseSigned := d.Get("response_signed").(bool)
	assertionSigned := d.Get("assertion_signed").(bool)

	preconfigName, isPreconfig := d.GetOkExists("preconfigured_app")

	if isPreconfig {
		app.Name = preconfigName.(string)
	} else {
		app.Name = d.Get("name").(string)

		reason := "Custom SAML applications must contain these fields"
		// Need to verify the fields that are now required since it is not preconfigured
		if err := conditionalRequire(d, customSamlAppRequiredFields, reason); err != nil {
			return app, err
		}

		if !responseSigned && !assertionSigned {
			return app, errors.New("custom SAML apps must either have response_signed or assertion_signed set to true")
		}
	}

	honorForce := d.Get("honor_force_authn").(bool)
	autoSubmit := d.Get("auto_submit_toolbar").(bool)
	hideMobile := d.Get("hide_ios").(bool)
	hideWeb := d.Get("hide_web").(bool)
	a11ySelfService := d.Get("accessibility_self_service").(bool)
	app.Settings = okta.NewSamlApplicationSettings()
	app.Visibility = &okta.ApplicationVisibility{
		AutoSubmitToolbar: &autoSubmit,
		Hide: &okta.ApplicationVisibilityHide{
			IOS: &hideMobile,
			Web: &hideWeb,
		},
	}
	app.Settings.SignOn = &okta.SamlApplicationSettingsSignOn{
		DefaultRelayState:     d.Get("default_relay_state").(string),
		SsoAcsUrl:             d.Get("sso_url").(string),
		Recipient:             d.Get("recipient").(string),
		Destination:           d.Get("destination").(string),
		Audience:              d.Get("audience").(string),
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
	app.Accessibility = &okta.ApplicationAccessibility{
		SelfService:      &a11ySelfService,
		ErrorRedirectUrl: d.Get("accessibility_error_redirect_url").(string),
		LoginRedirectUrl: d.Get("accessibility_login_redirect_url").(string),
	}

	return app, nil
}
