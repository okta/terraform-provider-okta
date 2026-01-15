package idaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/sdk"
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
	var (
		authenticator *sdk.Authenticator
		err           error
	)
	if id != "" {
		authenticator, _, err = getOktaClientFromMetadata(meta).Authenticator.GetAuthenticator(ctx, id)
	} else {
		authenticator, err = findAuthenticator(ctx, meta, name, key)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(authenticator.Id)
	_ = d.Set("key", authenticator.Key)
	_ = d.Set("name", authenticator.Name)
	_ = d.Set("status", authenticator.Status)
	_ = d.Set("type", authenticator.Type)
	if authenticator.Settings != nil {
		b, _ := json.Marshal(authenticator.Settings)
		_ = d.Set("settings", string(b))
	}
	if authenticator.Provider != nil {
		b, _ := json.Marshal(authenticator.Provider)
		dataMap := map[string]interface{}{}
		_ = json.Unmarshal([]byte(string(b)), &dataMap)
		b, _ = json.Marshal(dataMap)
		_ = d.Set("provider_json", string(b))

		_ = d.Set("provider_type", authenticator.Provider.Type)

		if authenticator.Type == "security_key" {
			_ = d.Set("provider_hostname", authenticator.Provider.Configuration.HostName)
			if authenticator.Provider.Configuration.AuthPortPtr != nil {
				_ = d.Set("provider_auth_port", authenticator.Provider.Configuration.AuthPortPtr)
			}
			_ = d.Set("provider_instance_id", authenticator.Provider.Configuration.InstanceId)
		}

		if authenticator.Provider.Type == "DUO" {
			_ = d.Set("provider_host", authenticator.Provider.Configuration.Host)
			_ = d.Set("provider_secret_key", authenticator.Provider.Configuration.SecretKey)
			_ = d.Set("provider_integration_key", authenticator.Provider.Configuration.IntegrationKey)
		}

		if authenticator.Provider.Configuration.UserNameTemplate != nil {
			_ = d.Set("provider_user_name_template", authenticator.Provider.Configuration.UserNameTemplate.Template)
		}
	}
	return nil
}

func findAuthenticator(ctx context.Context, meta interface{}, name, key string) (*sdk.Authenticator, error) {
	authenticators, _, err := getOktaClientFromMetadata(meta).Authenticator.ListAuthenticators(ctx)
	if err != nil {
		return nil, err
	}
	for _, authenticator := range authenticators {
		if key == "custom_app" {
			if authenticator.Name == name { // there can be more than 1 custom_app type authenticator, return nil in the end if we can't find by name.
				return authenticator, nil // TODO: update condition to include custom_otp as there can be more than 1 custom_otp type authenticator.
			}
		} else if key != "custom_otp" {
			if authenticator.Name == name {
				return authenticator, nil
			}
			if authenticator.Key == key {
				return authenticator, nil
			}
		} else {
			if authenticator.Name == name && authenticator.Key == key {
				return authenticator, nil
			} else {
				return nil, fmt.Errorf("authenticator with name '%s' and/or key '%s' does not exist", name, key)
			}
		}
	}
	if key != "" {
		return nil, fmt.Errorf("authenticator with key '%s' does not exist", key)
	}
	return nil, fmt.Errorf("authenticator with name '%s' does not exist", name) // authenticator names must be unique.
}
