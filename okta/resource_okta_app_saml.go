package okta

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

const (
	postBinding     = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
	redirectBinding = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
	saml11          = "1.1"
	saml20          = "2.0"
)

var (
	// Fields required if preconfigured_app is not provided
	customAppSamlRequiredFields = []string{
		"sso_url",
		"recipient",
		"destination",
		"audience",
		"subject_name_id_template",
		"subject_name_id_format",
		"signature_algorithm",
		"digest_algorithm",
		"authn_context_class_ref",
	}
	samlVersions = map[string]string{
		saml11: "SAML_1_1",
		saml20: "SAML_2_0",
	}
)

func isValidSkipArg(s string) bool {
	return s == "skip_users" || s == "skip_groups"
}

func resourceAppSaml() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSamlCreate,
		ReadContext:   resourceAppSamlRead,
		UpdateContext: resourceAppSamlUpdate,
		DeleteContext: resourceAppSamlDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
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
				Type:         schema.TypeString,
				Description:  "Certificate name. This modulates the rotation of keys. New name == new key.",
				Optional:     true,
				RequiredWith: []string{"key_years_valid"},
			},
			"key_id": {
				Type:        schema.TypeString,
				Description: "Certificate ID",
				Computed:    true,
			},
			"key_years_valid": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(2, 10),
				Description:      "Number of years the certificate is valid.",
			},
			"metadata": {
				Type:        schema.TypeString,
				Description: "SAML xml metadata payload",
				Computed:    true,
			},
			"metadata_url": {
				Type:        schema.TypeString,
				Description: "SAML xml metadata URL",
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
			"implicit_assignment": {
				Type:          schema.TypeBool,
				Optional:      true,
				Description:   "*Early Access Property*. Enable Federation Broker Mode.",
				ConflictsWith: []string{"groups", "users"},
			},
			"default_relay_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies a specific application resource in an IDP initiated SSO scenario.",
			},
			"sso_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Single Sign On URL",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"recipient": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The location where the app may present the SAML assertion",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"destination": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Identifies the location where the SAML response is intended to be sent inside of the SAML assertion",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"audience": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Audience Restriction",
			},
			"idp_issuer": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "SAML issuer ID",
				DiffSuppressFunc: appSamlDiffSuppressFunc,
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
				ValidateDiagFunc: elemInSlice([]string{
					"urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
					"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
					"urn:oasis:names:tc:SAML:1.1:nameid-format:x509SubjectName",
					"urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
					"urn:oasis:names:tc:SAML:2.0:nameid-format:transient",
				}),
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
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Signature algorithm used ot digitally sign the assertion and response",
				ValidateDiagFunc: elemInSlice([]string{"RSA_SHA256", "RSA_SHA1"}),
			},
			"digest_algorithm": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Determines the digest algorithm used to digitally sign the SAML assertion and response",
				ValidateDiagFunc: elemInSlice([]string{"SHA256", "SHA1"}),
			},
			"honor_force_authn": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Prompt user to re-authenticate if SP asks for it",
				Default:     false,
			},
			"authn_context_class_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies the SAML authentication context class for the assertion’s authentication statement",
			},
			"features": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "features to enable",
				Elem:        &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// always suppress diff since you can't currently configure provisioning features via the API
					return true
				},
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
			"inline_hook_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Saml Inline Hook setting",
			},
			"acs_endpoints": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of ACS endpoints for this SAML application",
			},
			"attribute_statements": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_type": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "Type of group attribute filter",
							ValidateDiagFunc: elemInSlice([]string{"STARTS_WITH", "EQUALS", "CONTAINS", "REGEX"}),
						},
						"filter_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter value to use",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The reference name of the attribute statement",
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified",
							ValidateDiagFunc: elemInSlice([]string{
								"urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified",
								"urn:oasis:names:tc:SAML:2.0:attrname-format:uri",
								"urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
							}),
							Description: "The name format of the attribute",
						},
						"type": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "EXPRESSION",
							ValidateDiagFunc: elemInSlice([]string{"GROUP", "EXPRESSION"}),
							Description:      "The type of attribute statements object",
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"single_logout_issuer": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The issuer of the Service Provider that generates the Single Logout request",
				RequiredWith: []string{"single_logout_url", "single_logout_certificate"},
			},
			"single_logout_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The location where the logout response is sent",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
				RequiredWith:     []string{"single_logout_issuer", "single_logout_certificate"},
			},
			"single_logout_certificate": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "x509 encoded certificate that the Service Provider uses to sign Single Logout requests",
				RequiredWith: []string{"single_logout_issuer", "single_logout_url"},
			},
			"saml_version": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          saml20,
				Description:      "SAML version for the app's sign-on mode",
				ValidateDiagFunc: elemInSlice([]string{saml20, saml11}),
			},
		}),
	}
}

func resourceAppSamlCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateAppSaml(d)
	if err != nil {
		return diag.FromErr(err)
	}
	app, err := buildSamlApp(d)
	if err != nil {
		return diag.Errorf("failed to create SAML application: %v", err)
	}
	activate := d.Get("status").(string) == statusActive
	params := &query.Params{Activate: &activate}
	_, _, err = getOktaClientFromMetadata(m).Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create SAML application: %v", err)
	}
	// Make sure to track in terraform prior to the creation of cert in case there is an error.
	d.SetId(app.Id)
	// When the implicit_assignment is turned on, calls to the user/group assignments will error with a bad request
	// So Skip setting assignments while this is on
	if !d.Get("implicit_assignment").(bool) {
		err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
		if err != nil {
			return diag.Errorf("failed to handle groups and users for SAML application: %v", err)
		}
	}
	err = tryCreateCertificate(ctx, d, m, app.Id)
	if err != nil {
		return diag.Errorf("failed to create new certificate for SAML application: %v", err)
	}
	err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to handle groups and users for SAML application: %v", err)
	}
	err = handleAppLogo(ctx, d, m, app.Id, app.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for SAML application: %v", err)
	}
	return resourceAppSamlRead(ctx, d, m)
}

func resourceAppSamlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := okta.NewSamlApplication()
	err := fetchApp(ctx, d, m, app)
	if err != nil {
		return diag.Errorf("failed to get SAML application: %v", err)
	}
	if app.Id == "" {
		d.SetId("")
		return nil
	}
	if app.Settings != nil {
		if app.Settings.SignOn != nil {
			err = setSamlSettings(d, app.Settings.SignOn)
			if err != nil {
				return diag.Errorf("failed to set SAML sign-on settings: %v", err)
			}
		}
		err = setAppSettings(d, app.Settings.App)
		if err != nil {
			return diag.Errorf("failed to set SAML app settings: %v", err)
		}
	}
	_ = d.Set("features", convertStringSliceToSetNullable(app.Features))
	_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
	_ = d.Set("user_name_template_push_status", app.Credentials.UserNameTemplate.PushStatus)
	_ = d.Set("preconfigured_app", app.Name)
	_ = d.Set("logo_url", linksValue(app.Links, "logo", "href"))
	if app.Settings.ImplicitAssignment != nil {
		_ = d.Set("implicit_assignment", *app.Settings.ImplicitAssignment)
	} else {
		_ = d.Set("implicit_assignment", false)
	}
	if app.Credentials.Signing.Kid != "" && app.Status != statusInactive {
		keyID := app.Credentials.Signing.Kid
		_ = d.Set("key_id", keyID)
		keyMetadata, metadataRoot, err := getSupplementFromMetadata(m).GetSAMLMetadata(ctx, d.Id(), keyID)
		if err != nil {
			return diag.Errorf("failed to get app's SAML metadata: %v", err)
		}
		var q string
		if keyID != "" {
			q = fmt.Sprintf("?kid=%s", keyID)
		}
		_ = d.Set("metadata", string(keyMetadata))
		_ = d.Set("metadata_url", fmt.Sprintf("%s/api/v1/apps/%s/sso/saml/metadata%s",
			getOktaClientFromMetadata(m).GetConfig().Okta.Client.OrgUrl, d.Id(), q))
		desc := metadataRoot.IDPSSODescriptors[0]
		syncSamlEndpointBinding(d, desc.SingleSignOnServices)
		uri := metadataRoot.EntityID
		key := getExternalID(uri, app.Settings.SignOn.IdpIssuer)
		_ = d.Set("entity_url", uri)
		_ = d.Set("entity_key", key)
		_ = d.Set("certificate", desc.KeyDescriptors[0].KeyInfo.Certificate)
	}
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	if app.SignOnMode == "SAML_1_1" {
		_ = d.Set("saml_version", saml11)
	} else {
		_ = d.Set("saml_version", saml20)
	}
	// When the implicit_assignment is turned on, calls to the user/group assignments will error with a bad request
	// So Skip setting assignments while this is on
	if !d.Get("implicit_assignment").(bool) {
		if err = syncGroupsAndUsers(ctx, app.Id, d, m); err != nil {
			return diag.Errorf("failed to sync groups and users for OAuth application: %v", err)
		}
	}
	return nil
}

func resourceAppSamlUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateAppSaml(d)
	if err != nil {
		return diag.FromErr(err)
	}
	client := getOktaClientFromMetadata(m)
	app, err := buildSamlApp(d)
	if err != nil {
		return diag.Errorf("failed to create SAML application: %v", err)
	}
	err = setAppStatus(ctx, d, client, app.Status)
	if err != nil {
		return diag.Errorf("failed to set SAML application status: %v", err)
	}
	_, _, err = client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update SAML application: %v", err)
	}
	if d.HasChange("key_name") {
		err = tryCreateCertificate(ctx, d, m, app.Id)
		if err != nil {
			return diag.Errorf("failed to create new certificate for SAML application: %v", err)
		}
	}
	// When the implicit_assignment is turned on, calls to the user/group assignments will error with a bad request
	// So Skip setting assignments while this is on
	if !d.Get("implicit_assignment").(bool) {
		err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
		if err != nil {
			return diag.Errorf("failed to handle groups and users for OAuth application: %v", err)
		}
	}
	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, m, app.Id, app.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for SAML application: %v", err)
		}
	}
	isStatusChaged := d.HasChange("status")
	if isStatusChaged {
		s := d.Get("status").(string)
		if s == "ACTIVE" {
			// activate
		} else {
			// deactivate
		}
	}
	return resourceAppSamlRead(ctx, d, m)
}

func resourceAppSamlDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to delete SAML application: %v", err)
	}
	return nil
}

func buildSamlApp(d *schema.ResourceData) (*okta.SamlApplication, error) {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewSamlApplication()
	app.SignOnMode = samlVersions[d.Get("saml_version").(string)]
	app.Label = d.Get("label").(string)
	responseSigned := d.Get("response_signed").(bool)
	assertionSigned := d.Get("assertion_signed").(bool)

	preConfigName, ok := d.GetOk("preconfigured_app")
	if ok {
		app.Name = preConfigName.(string)
	} else {
		app.Name = d.Get("name").(string)

		reason := "Custom SAML applications must contain these fields"
		// Need to verify the fields that are now required since it is not preconfigured
		if err := conditionalRequire(d, customAppSamlRequiredFields, reason); err != nil {
			return app, err
		}

		if !responseSigned && !assertionSigned {
			return app, errors.New("custom SAML apps must either have response_signed or assertion_signed set to true")
		}
	}

	honorForce := d.Get("honor_force_authn").(bool)
	app.Settings = &okta.SamlApplicationSettings{
		ImplicitAssignment: boolPtr(d.Get("implicit_assignment").(bool)),
		Notes:              buildAppNotes(d),
	}
	app.Visibility = buildAppVisibility(d)
	app.Accessibility = buildAppAccessibility(d)
	app.Settings.App = buildAppSettings(d)
	// Note: You can't currently configure provisioning features via the API. Use the administrator UI.
	// app.Features = convertInterfaceToStringSet(d.Get("features"))
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
		Slo:                   &okta.SingleLogout{Enabled: boolPtr(false)},
	}
	sli := d.Get("single_logout_issuer").(string)
	if sli != "" {
		app.Settings.SignOn.Slo = &okta.SingleLogout{
			Enabled:   boolPtr(true),
			Issuer:    sli,
			LogoutUrl: d.Get("single_logout_url").(string),
		}
		app.Settings.SignOn.SpCertificate = &okta.SpCertificate{
			X5c: []string{d.Get("single_logout_certificate").(string)},
		}
	}
	app.Credentials = &okta.ApplicationCredentials{
		UserNameTemplate: buildUserNameTemplate(d),
	}

	// Assumes that sso url is already part of the acs endpoints as part of the desired state.
	acsEndpoints := convertInterfaceToStringSet(d.Get("acs_endpoints"))

	// If there are acs endpoints, implies this flag should be true.
	allowMultipleAcsEndpoints := false
	if len(acsEndpoints) > 0 {
		acsEndpointsObj := make([]*okta.AcsEndpoint, len(acsEndpoints))
		for i := range acsEndpoints {
			acsEndpointsObj[i] = &okta.AcsEndpoint{
				Index: int64(i),
				Url:   acsEndpoints[i],
			}
		}
		allowMultipleAcsEndpoints = true
		app.Settings.SignOn.AcsEndpoints = acsEndpointsObj
	}
	app.Settings.SignOn.AllowMultipleAcsEndpoints = &allowMultipleAcsEndpoints

	statements := d.Get("attribute_statements").([]interface{})
	if len(statements) > 0 {
		samlAttr := make([]*okta.SamlAttributeStatement, len(statements))
		for i := range statements {
			samlAttr[i] = &okta.SamlAttributeStatement{
				FilterType:  d.Get(fmt.Sprintf("attribute_statements.%d.filter_type", i)).(string),
				FilterValue: d.Get(fmt.Sprintf("attribute_statements.%d.filter_value", i)).(string),
				Name:        d.Get(fmt.Sprintf("attribute_statements.%d.name", i)).(string),
				Namespace:   d.Get(fmt.Sprintf("attribute_statements.%d.namespace", i)).(string),
				Type:        d.Get(fmt.Sprintf("attribute_statements.%d.type", i)).(string),
				Values:      convertInterfaceToStringArr(d.Get(fmt.Sprintf("attribute_statements.%d.values", i))),
			}
		}
		app.Settings.SignOn.AttributeStatements = samlAttr
	} else {
		app.Settings.SignOn.AttributeStatements = []*okta.SamlAttributeStatement{}
	}

	if id, ok := d.GetOk("key_id"); ok {
		app.Credentials.Signing = &okta.ApplicationCredentialsSigning{
			Kid: id.(string),
		}
	}

	if id, ok := d.GetOk("inline_hook_id"); ok {
		app.Settings.SignOn.InlineHooks = []*okta.SignOnInlineHook{{Id: id.(string)}}
	}

	return app, nil
}

// Keep in mind that at the time of writing this the official SDK did not support generating certs.
func generateCertificate(ctx context.Context, d *schema.ResourceData, m interface{}, appID string) (*okta.JsonWebKey, error) {
	requestExecutor := getRequestExecutor(m)
	years := d.Get("key_years_valid").(int)
	url := fmt.Sprintf("/api/v1/apps/%s/credentials/keys/generate?validityYears=%d", appID, years)
	req, err := requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	var key *okta.JsonWebKey
	_, err = requestExecutor.Do(ctx, req, &key)
	return key, err
}

func tryCreateCertificate(ctx context.Context, d *schema.ResourceData, m interface{}, appID string) error {
	if _, ok := d.GetOk("key_name"); ok {
		key, err := generateCertificate(ctx, d, m, appID)
		if err != nil {
			return err
		}

		// Set ID and the read done at the end of update and create will do the GET on metadata
		_ = d.Set("key_id", key.Kid)
	}
	return nil
}

func validateAppSaml(d *schema.ResourceData) error {
	jwks, ok := d.GetOk("attribute_statements")
	if !ok {
		return nil
	}
	for i := range jwks.([]interface{}) {
		objType := d.Get(fmt.Sprintf("attribute_statements.%d.type", i)).(string)
		if (d.Get(fmt.Sprintf("attribute_statements.%d.filter_type", i)).(string) != "" ||
			d.Get(fmt.Sprintf("attribute_statements.%d.filter_value", i)).(string) != "") &&
			objType != "GROUP" {
			return errors.New("invalid 'attribute_statements': when setting 'filter_value' or 'filter_type', value of 'type' should be set to 'GROUP'")
		}
		if objType == "GROUP" &&
			len(convertInterfaceToStringArrNullable(d.Get(fmt.Sprintf("attribute_statements.%d.values", i)))) > 0 {
			return errors.New("invalid 'attribute_statements': when setting 'values', 'type' should be set to 'EXPRESSION'")
		}
	}
	return nil
}
