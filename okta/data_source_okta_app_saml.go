package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppSaml() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppSamlRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label", "label_prefix"},
			},
			"label": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label_prefix"},
			},
			"label_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label"},
			},
			"active_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Search only ACTIVE applications.",
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
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
				ValidateDiagFunc: stringInSlice(
					[]string{
						"urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
						"urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
						"urn:oasis:names:tc:SAML:1.1:nameid-format:x509SubjectName",
						"urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
						"urn:oasis:names:tc:SAML:2.0:nameid-format:transient",
					},
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
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Signature algorithm used ot digitally sign the assertion and response",
				ValidateDiagFunc: stringInSlice([]string{"RSA_SHA256", "RSA_SHA1"}),
			},
			"digest_algorithm": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Determines the digest algorithm used to digitally sign the SAML assertion and response",
				ValidateDiagFunc: stringInSlice([]string{"SHA256", "SHA1"}),
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
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Custom error page URL",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"accessibility_login_redirect_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Custom login page URL",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"features": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "features to enable",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"user_name_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "${source.login}",
				Description: "Username template",
			},
			"user_name_template_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username template suffix",
			},
			"user_name_template_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "BUILT_IN",
				Description:      "Username template type",
				ValidateDiagFunc: stringInSlice([]string{"NONE", "CUSTOM", "BUILT_IN"}),
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
					Schema: attributeStatements,
				},
			},
		},
	}
}

func dataSourceAppSamlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	filters, err := getAppFilters(d)
	if err != nil {
		return diag.Errorf("failed to list SAML apps: error getting filters: %v", err)
	}
	appList, err := listSamlApps(ctx, m.(*Config), filters)
	if err != nil {
		return diag.Errorf("failed to list SAML apps: error getting SAML apps: %v", err)
	}
	if len(appList) < 1 {
		return diag.Errorf("no SAML applications found with provided filter: %+v", filters)
	} else if len(appList) > 1 {
		logger(m).Info("found multiple applications with the criteria supplied, using the first one, sorted by creation date")
	}
	app := appList[0]
	d.SetId(app.Id)
	_ = d.Set("label", app.Label)
	_ = d.Set("name", app.Name)
	_ = d.Set("status", app.Status)
	_ = d.Set("key_id", app.Credentials.Signing.Kid)
	if app.Settings != nil && app.Settings.SignOn != nil {
		err = syncSamlSettings(d, app.Settings)
		if err != nil {
			return diag.Errorf("failed to read SAML app: error setting SAML app settings: %v", err)
		}
	}
	_ = d.Set("features", convertStringSetToInterface(app.Features))
	_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
	return nil
}
