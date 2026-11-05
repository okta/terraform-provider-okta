package idaas

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func dataSourceAuthenticator() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthenticatorRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key", "name"},
				Description:   "ID of the authenticator.",
			},
			"key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "name"},
				Description:   "A human-readable string that identifies the authenticator.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "key"},
				Description:   "Name of the authenticator.",
			},
			"settings": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Authenticator settings in JSON format",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Authenticator.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the authenticator",
			},
			"provider_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Authenticator Provider in JSON format",
			},
			"provider_auth_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The RADIUS server port (for example 1812). This is defined when the On-Prem RADIUS server is configured",
			},
			"provider_hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server host name or IP address",
			},
			"provider_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "(Specific to `security_key`) App Instance ID.",
			},
			"provider_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provider type.",
			},
			"provider_user_name_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username template expected by the provider.",
			},
		},
		Description: "Get an authenticator by key, name of ID.",
	}
}

func dataSourceAuthenticatorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return datasourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticator)
	}

	id := d.Get("id").(string)
	name := d.Get("name").(string)
	key := d.Get("key").(string)
	if id == "" && name == "" && key == "" {
		return diag.Errorf("config must provide either 'id', 'name' or 'key' to retrieve the authenticator")
	}

	v6Client := getOktaV6ClientFromMetadata(meta)
	var err error
	var authenticator *v6okta.AuthenticatorBase
	if id != "" {
		authenticator, _, err = v6Client.AuthenticatorAPI.GetAuthenticator(ctx, id).Execute()
	} else {
		authenticator, err = findAuthenticatorV6(ctx, v6Client, name, key)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if authenticator == nil {
		if key != "" {
			return diag.Errorf("authenticator with key '%s' does not exist", key)
		}
		return diag.Errorf("authenticator with name '%s' does not exist", name)
	}

	d.SetId(authenticator.GetId())
	establishAuthenticatorDataSourceV6(authenticator, d)
	return nil
}

// establishAuthenticatorDataSourceV6 mirrors establishAuthenticatorV6, but
// always populates provider_json (this data source's provider_json attribute
// is Computed, unlike the resource's, so there's no risk of a perpetual diff
// from always setting it).
func establishAuthenticatorDataSourceV6(authenticator *v6okta.AuthenticatorBase, d *schema.ResourceData) {
	_ = d.Set("key", authenticator.GetKey())
	_ = d.Set("name", authenticator.GetName())
	_ = d.Set("status", authenticator.GetStatus())
	_ = d.Set("type", authenticator.GetType())

	if raw, ok := authenticator.AdditionalProperties["settings"]; ok && raw != nil {
		if s, ok := marshalAuthenticatorSettings(raw); ok {
			_ = d.Set("settings", s)
		}
	}

	raw, ok := authenticator.AdditionalProperties["provider"]
	if !ok || raw == nil {
		return
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return
	}
	_ = d.Set("provider_json", string(b))

	var provider authenticatorProviderV6
	if json.Unmarshal(b, &provider) != nil {
		return
	}
	_ = d.Set("provider_type", provider.Type)
	if provider.Configuration == nil {
		return
	}

	if authenticator.GetType() == "security_key" {
		_ = d.Set("provider_hostname", provider.Configuration.HostName)
		if provider.Configuration.AuthPort != nil {
			_ = d.Set("provider_auth_port", int(*provider.Configuration.AuthPort))
		}
		_ = d.Set("provider_instance_id", provider.Configuration.InstanceId)
	}

	if provider.Type == "DUO" {
		_ = d.Set("provider_host", provider.Configuration.Host)
		_ = d.Set("provider_secret_key", provider.Configuration.SecretKey)
		_ = d.Set("provider_integration_key", provider.Configuration.IntegrationKey)
	}

	if provider.Configuration.UserNameTemplate != nil {
		_ = d.Set("provider_user_name_template", provider.Configuration.UserNameTemplate.Template)
	}
}
