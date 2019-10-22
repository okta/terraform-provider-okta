package okta

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/crewjam/saml"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

const (
	postBinding     = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
	redirectBinding = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
)

// Fields required if preconfigured_app is not provided
var customappSamlRequiredFields = []string{
	"sso_url",
	"recipient",
	"destination",
	"audience",
	"subject_name_id_template",
	"subject_name_id_format",
	"signature_algorithm",
	"digest_algorithm",
	"honor_force_authn",
	"authn_context_class_ref",
}

func resourceAppSaml() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppSamlCreate,
		Read:   resourceAppSamlRead,
		Update: resourceAppSamlUpdate,
		Delete: resourceAppSamlDelete,
		Exists: resourceAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// For those familiar with Terraform schemas be sure to check the base application schema and/or
		// the examples in the documentation
		Schema: buildAppSchema(map[string]*schema.Schema{
			"preconfigured_app": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of preexisting SAML application. For instance 'slack'",
				ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"key_name": {
				Type:        schema.TypeString,
				Description: "Certificate name. This modulates the rotation of keys. New name == new key.",
				Optional:    true,
			},
			"key_id": {
				Type:        schema.TypeString,
				Description: "Certificate ID",
				Computed:    true,
			},
			"key_years_valid": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Number of years the certificate is valid.",
			},
			"metadata": {
				Type:        schema.TypeString,
				Description: "SAML xml metadata payload",
				Computed:    true,
			},
			"certificate": {
				Type:        schema.TypeString,
				Description: "cert from SAML XML metadata payload",
				Computed:    true,
			},
			"http_post_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Post location from the SAML metadata.",
			},
			"http_redirect_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect location from the SAML metadata.",
			},
			"entity_key": {
				Type:        schema.TypeString,
				Description: "Entity ID, the ID portion of the entity_url",
				Computed:    true,
			},
			"entity_url": {
				Type:        schema.TypeString,
				Description: "Entity URL for instance http://www.okta.com/exk1fcia6d6EMsf331d8",
				Computed:    true,
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
				Default:     false,
				Description: "Do not display application icon on mobile app",
			},
			"hide_web": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not display application icon to users",
			},
			"default_relay_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies a specific application resource in an IDP initiated SSO scenario.",
			},
			"sso_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Single Sign On URL",
				ValidateFunc: validateIsURL,
			},
			"recipient": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The location where the app may present the SAML assertion",
				ValidateFunc: validateIsURL,
			},
			"destination": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Identifies the location where the SAML response is intended to be sent inside of the SAML assertion",
				ValidateFunc: validateIsURL,
			},
			"audience": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Audience Restriction",
			},
			"idp_issuer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SAML issuer ID",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Conditional default
					return new == "" && old == "http://www.okta.com/${org.externalKey}"
				},
			},
			"sp_issuer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SAML SP issuer ID",
			},
			"subject_name_id_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template for app user's username when a user is assigned to the app",
			},
			"subject_name_id_format": {
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
			"response_signed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether the SAML auth response message is digitally signed",
			},
			"request_compressed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Denotes whether the request is compressed or not.",
			},
			"assertion_signed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether the SAML assertion is digitally signed",
			},
			"signature_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Signature algorithm used ot digitally sign the assertion and response",
				ValidateFunc: validation.StringInSlice([]string{"RSA_SHA256", "RSA_SHA1"}, false),
			},
			"digest_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Determines the digest algorithm used to digitally sign the SAML assertion and response",
				ValidateFunc: validation.StringInSlice([]string{"SHA256", "SHA1"}, false),
			},
			"honor_force_authn": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Prompt user to re-authenticate if SP asks for it",
			},
			"authn_context_class_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies the SAML authentication context class for the assertion’s authentication statement",
			},
			"accessibility_self_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable self service",
			},
			"accessibility_error_redirect_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Custom error page URL",
				ValidateFunc: validateIsURL,
			},
			"accessibility_login_redirect_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Custom login page URL",
				ValidateFunc: validateIsURL,
			},
			"features": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "features to enable",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"user_name_template": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "${source.login}",
				Description: "Username template",
			},
			"user_name_template_suffix": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username template suffix",
			},
			"user_name_template_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "BUILT_IN",
				Description:  "Username template type",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "CUSTOM", "BUILT_IN"}, false),
			},
			"app_settings_json": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Application settings in JSON format",
				ValidateFunc: validateDataJSON,
				StateFunc:    normalizeDataJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"attribute_statements": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of group attribute filter",
						},
						"filter_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter value to use",
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified",
							ValidateFunc: validation.StringInSlice([]string{
								"urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified",
								"urn:oasis:names:tc:SAML:2.0:attrname-format:uri",
								"urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
							}, false),
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "EXPRESSION",
							ValidateFunc: validation.StringInSlice([]string{
								"EXPRESSION",
								"GROUP",
							}, false),
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		}),
	}
}

func resourceAppSamlCreate(d *schema.ResourceData, m interface{}) error {
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

	// Make sure to track in terraform prior to the creation of cert in case there is an error.
	d.SetId(app.Id)

	err = tryCreateCertificate(d, m, app.Id)
	if err != nil {
		return err
	}
	err = handleAppGroupsAndUsers(app.Id, d, m)

	if err != nil {
		return err
	}

	return resourceAppSamlRead(d, m)
}

func resourceAppSamlRead(d *schema.ResourceData, m interface{}) error {
	app := okta.NewSamlApplication()
	err := fetchApp(d, m, app)

	if app == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	if app.Settings != nil && app.Settings.SignOn != nil {
		syncSamlSettings(d, app.Settings)
	}

	d.Set("features", convertStringSetToInterface(app.Features))
	d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
	d.Set("preconfigured_app", app.Name)
	err = setAppSettings(d, app.Settings.App)
	if err != nil {
		return err
	}

	if app.Credentials.Signing.Kid != "" && app.Status != "INACTIVE" {
		keyID := app.Credentials.Signing.Kid
		d.Set("key_id", keyID)
		keyMetadata, err := getMetadata(d, m, keyID)
		if err != nil {
			return err
		}
		d.Set("metadata", string(keyMetadata))

		metadataRoot := &saml.EntityDescriptor{}
		err = xml.Unmarshal(keyMetadata, metadataRoot)
		if err != nil {
			return fmt.Errorf("Could not parse SAML app metadata, error: %s", err)
		}
		desc := metadataRoot.IDPSSODescriptors[0]
		syncSamlEndpointBinding(d, desc.SingleSignOnServices)
		uri := metadataRoot.EntityID
		key := getExternalID(uri, app.Settings.SignOn.IdpIssuer)
		d.Set("entity_url", uri)
		d.Set("entity_key", key)
		d.Set("certificate", desc.KeyDescriptors[0].KeyInfo.Certificate)
	}

	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility)

	return syncGroupsAndUsers(app.Id, d, m)
}

func resourceAppSamlUpdate(d *schema.ResourceData, m interface{}) error {
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

	if d.HasChange("key_name") {
		err = tryCreateCertificate(d, m, app.Id)
		if err != nil {
			return err
		}
	}
	err = handleAppGroupsAndUsers(app.Id, d, m)

	if err != nil {
		return err
	}

	return resourceAppSamlRead(d, m)
}

func resourceAppSamlDelete(d *schema.ResourceData, m interface{}) error {
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
	responseSigned := d.Get("response_signed").(bool)
	assertionSigned := d.Get("assertion_signed").(bool)

	preconfigName, isPreconfig := d.GetOkExists("preconfigured_app")

	if isPreconfig {
		app.Name = preconfigName.(string)
	} else {
		app.Name = d.Get("name").(string)

		reason := "Custom SAML applications must contain these fields"
		// Need to verify the fields that are now required since it is not preconfigured
		if err := conditionalRequire(d, customappSamlRequiredFields, reason); err != nil {
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
	if appSettings, ok := d.GetOk("app_settings_json"); ok {
		payload := map[string]interface{}{}
		json.Unmarshal([]byte(appSettings.(string)), &payload)
		settings := okta.ApplicationSettingsApplication(payload)
		app.Settings.App = &settings
	}
	app.Features = convertInterfaceToStringSet(d.Get("features"))
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
	app.Credentials = &okta.ApplicationCredentials{
		UserNameTemplate: &okta.ApplicationCredentialsUsernameTemplate{
			Template: d.Get("user_name_template").(string),
			Type:     d.Get("user_name_template_type").(string),
			Suffix:   d.Get("user_name_template_suffix").(string),
		},
	}
	app.Accessibility = &okta.ApplicationAccessibility{
		SelfService:      &a11ySelfService,
		ErrorRedirectUrl: d.Get("accessibility_error_redirect_url").(string),
		LoginRedirectUrl: d.Get("accessibility_login_redirect_url").(string),
	}
	statements := d.Get("attribute_statements").([]interface{})
	if len(statements) > 0 {
		samlAttr := make([]*okta.SamlAttributeStatement, len(statements))
		for i := range statements {
			samlAttr[i] = &okta.SamlAttributeStatement{
				Name:        d.Get(fmt.Sprintf("attribute_statements.%d.name", i)).(string),
				Namespace:   d.Get(fmt.Sprintf("attribute_statements.%d.namespace", i)).(string),
				Type:        d.Get(fmt.Sprintf("attribute_statements.%d.type", i)).(string),
				Values:      convertInterfaceToStringArr(d.Get(fmt.Sprintf("attribute_statements.%d.values", i))),
				FilterType:  d.Get(fmt.Sprintf("attribute_statements.%d.filter_type", i)).(string),
				FilterValue: d.Get(fmt.Sprintf("attribute_statements.%d.filter_value", i)).(string),
			}
		}
		app.Settings.SignOn.AttributeStatements = samlAttr
	}

	if id, ok := d.GetOk("key_id"); ok {
		app.Credentials = &okta.ApplicationCredentials{
			Signing: &okta.ApplicationCredentialsSigning{
				Kid: id.(string),
			},
		}
	}

	return app, nil
}

func getCertificate(d *schema.ResourceData, m interface{}) (*okta.JsonWebKey, error) {
	client := getOktaClientFromMetadata(m)
	keyId := d.Get("key.id").(string)
	key, resp, err := client.Application.GetApplicationKey(d.Id(), keyId)
	if resp.StatusCode == 404 {
		return nil, nil
	}

	return key, err
}

func getMetadata(d *schema.ResourceData, m interface{}, keyID string) ([]byte, error) {
	key, _, err := getSupplementFromMetadata(m).GetSAMLMetdata(d.Id(), keyID)
	return key, err
}

// Keep in mind that at the time of writing this the official SDK did not support generating certs.
func generateCertificate(d *schema.ResourceData, m interface{}, appID string) (*okta.JsonWebKey, error) {
	requestExecutor := getRequestExecutor(m)
	years := d.Get("key_years_valid").(int)
	url := fmt.Sprintf("/api/v1/apps/%s/credentials/keys/generate?validityYears=%d", appID, years)
	req, err := requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	var key *okta.JsonWebKey

	_, err = requestExecutor.Do(req, &key)

	return key, err
}

func tryCreateCertificate(d *schema.ResourceData, m interface{}, appID string) error {
	if _, ok := d.GetOk("key_name"); ok {
		key, err := generateCertificate(d, m, appID)
		if err != nil {
			return err
		}

		// Set ID and the read done at the end of update and create will do the GET on metadata
		d.Set("key_id", key.Kid)
	}

	return nil
}
