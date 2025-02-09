package okta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceAppSaml() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppSamlRead,
		Schema: buildSchema(skipUsersAndGroupsSchema, map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label", "label_prefix"},
				Description:   "Id of application to retrieve, conflicts with label and label_prefix.",
			},
			"label": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label_prefix"},
				Description: `The label of the app to retrieve, conflicts with label_prefix and id. Label 
				uses the ?q=<label> query parameter exposed by Okta's API. It should be noted that at this time 
				this searches both name and label. This is used to avoid paginating through all applications.`,
			},
			"label_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label"},
				Description: `Label prefix of the app to retrieve, conflicts with label and id. This will tell the
				provider to do a starts with query as opposed to an equals query.`,
			},
			"active_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Search only ACTIVE applications.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of application.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of application.",
			},
			"key_id": {
				Type:        schema.TypeString,
				Description: "Certificate ID",
				Computed:    true,
			},
			"auto_submit_toolbar": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Display auto submit toolbar",
			},
			"hide_ios": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Do not display application icon on mobile app",
			},
			"hide_web": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Do not display application icon to users",
			},
			"default_relay_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifies a specific application resource in an IDP initiated SSO scenario.",
			},
			"sso_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Single Sign On URL",
			},
			"recipient": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location where the app may present the SAML assertion",
			},
			"destination": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifies the location where the SAML response is intended to be sent inside of the SAML assertion",
			},
			"audience": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Audience Restriction",
			},
			"idp_issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SAML issuer ID",
			},
			"sp_issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SAML SP issuer ID",
			},
			"subject_name_id_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template for app user's username when a user is assigned to the app",
			},
			"subject_name_id_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifies the SAML processing rules.",
			},
			"response_signed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether the SAML auth response message is digitally signed",
			},
			"request_compressed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Denotes whether the request is compressed or not.",
			},
			"assertion_signed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether the SAML assertion is digitally signed",
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Signature algorithm used to digitally sign the assertion and response",
			},
			"digest_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Determines the digest algorithm used to digitally sign the SAML assertion and response",
			},
			"honor_force_authn": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Prompt user to re-authenticate if SP asks for it",
			},
			"authn_context_class_ref": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifies the SAML authentication context class for the assertion’s authentication statement",
			},
			"accessibility_self_service": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable self service",
			},
			"accessibility_error_redirect_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Custom error page URL",
			},
			"accessibility_login_redirect_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Custom login page URL",
			},
			"features": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "features to enable",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"user_name_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username template",
			},
			"user_name_template_suffix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username template suffix",
			},
			"user_name_template_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username template type",
			},
			"user_name_template_push_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Push username on update",
			},
			"app_settings_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Application settings in JSON format",
			},
			"acs_endpoints": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "List of ACS endpoints for this SAML application",
			},
			"attribute_statements": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of group attribute filter",
						},
						"filter_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Filter value to use",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reference name of the attribute statement",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name format of the attribute",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of attribute statements object",
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"single_logout_issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The issuer of the Service Provider that generates the Single Logout request",
			},
			"single_logout_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The location where the logout response is sent",
			},
			"single_logout_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "x509 encoded certificate that the Service Provider uses to sign Single Logout requests",
			},
			"links": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Discoverable resources related to the app",
			},
			"inline_hook_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Saml Inline Hook setting",
			},
			"groups": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Groups associated with the application",
				Deprecated:  "The `groups` field is now deprecated for the data source `okta_app_saml`, please replace all uses of this with: `okta_app_group_assignments`",
			},
			"users": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Users associated with the application",
				Deprecated:  "The `users` field is now deprecated for the data source `okta_app_saml`, please replace all uses of this with: `okta_app_user_assignments`",
			},
			"saml_signed_request_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "SAML Signed Request enabled",
			},
		}),
		Description: "Get a SAML application from Okta.",
	}
}

func dataSourceAppSamlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	filters, err := getAppFilters(d)
	if err != nil {
		return diag.Errorf("invalid SAML app filters: %v", err)
	}
	var app *sdk.SamlApplication
	if filters.ID != "" {
		respApp, _, err := getOktaClientFromMetadata(m).Application.GetApplication(ctx, filters.ID, sdk.NewSamlApplication(), nil)
		if err != nil {
			return diag.Errorf("failed get app by ID: %v", err)
		}
		app = respApp.(*sdk.SamlApplication)
	} else {
		re := getOktaClientFromMetadata(m).GetRequestExecutor()
		qp := &query.Params{Filter: filters.Status, Q: filters.getQ()}
		req, err := re.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/apps%s", qp.String()), nil)
		if err != nil {
			return diag.Errorf("failed to list SAML apps: %v", err)
		}
		var appList []*sdk.SamlApplication
		_, err = re.Do(ctx, req, &appList)
		if err != nil {
			return diag.Errorf("failed to list SAML apps: %v", err)
		}
		if len(appList) < 1 {
			return diag.Errorf("no SAML application found with provided filter: %s", filters)
		}

		if filters.Label != "" {
			foundMatch := false
			for _, appItx := range appList {
				if appItx.Label == filters.Label {
					app = appItx
					foundMatch = true
					break
				}
			}
			if !foundMatch {
				return diag.Errorf("no SAML application found with the provided label: %s", filters.Label)
			}
		} else {
			logger(m).Info("found multiple SAML applications with the criteria supplied, using the first one, sorted by creation date")
			app = appList[0]
		}
	}
	d.SetId(app.Id)
	_ = d.Set("label", app.Label)
	_ = d.Set("name", app.Name)
	_ = d.Set("status", app.Status)
	_ = d.Set("key_id", app.Credentials.Signing.Kid)
	if app.Settings != nil {
		if app.Settings.SignOn != nil {
			err = setSamlSettings(d, app.Settings.SignOn)
			if err != nil {
				return diag.Errorf("failed to read SAML app: error setting SAML sign-on settings: %v", err)
			}
		}
		err = setAppSettings(d, app.Settings.App)
		if err != nil {
			return diag.Errorf("failed to read SAML app: failed to set SAML app settings: %v", err)
		}
	}
	_ = d.Set("features", convertStringSliceToSetNullable(app.Features))
	_ = d.Set("user_name_template", app.Credentials.UserNameTemplate.Template)
	_ = d.Set("user_name_template_type", app.Credentials.UserNameTemplate.Type)
	_ = d.Set("user_name_template_suffix", app.Credentials.UserNameTemplate.Suffix)
	_ = d.Set("user_name_template_push_status", app.Credentials.UserNameTemplate.PushStatus)
	p, _ := json.Marshal(app.Links)
	_ = d.Set("links", string(p))
	return nil
}
