package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAppSaml() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppSamlRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label", "label_prefix"},
			},
			"label": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label_prefix"},
			},
			"label_prefix": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label"},
			},
			"active_only": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Search only ACTIVE applications.",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_id": {
				Type:        schema.TypeString,
				Description: "Certificate ID",
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
		},
	}
}

func dataSourceAppSamlRead(d *schema.ResourceData, m interface{}) error {
	filters, err := getAppFilters(d)
	if err != nil {
		return err
	}

	appList, err := listSamlApps(m.(*Config), filters)
	if err != nil {
		return err
	}
	if len(appList) < 1 {
		return fmt.Errorf("No application found with provided filter: %s", filters)
	} else if len(appList) > 1 {
		fmt.Println("Found multiple applications with the criteria supplied, using the first one, sorted by creation date.")
	}
	app := appList[0]
	d.SetId(app.Id)
	d.Set("label", app.Label)
	d.Set("name", app.Name)
	d.Set("status", app.Status)
	d.Set("key_id", app.Credentials.Signing.Kid)

	if app.Settings != nil && app.Settings.SignOn != nil {
		syncSamlSettings(d, app.Settings)
	}

	d.Set("features", convertStringSetToInterface(app.Features))
	d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)

	return nil
}
