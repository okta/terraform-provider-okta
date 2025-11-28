package idaas

import (
	"context"
	"encoding/json"
	"fmt"

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

	client := getOktaV6ClientFromMetadata(meta)

	id := d.Get("id").(string)
	name := d.Get("name").(string)
	key := d.Get("key").(string)
	if id == "" && name == "" && key == "" {
		return diag.Errorf("config must provide either 'id', 'name' or 'key' to retrieve the authenticator")
	}
	var (
		authenticator *v6okta.AuthenticatorBase
		err           error
	)
	if id != "" {
		authenticator, _, err = client.AuthenticatorAPI.GetAuthenticator(ctx, id).Execute()
	} else {
		authenticator, err = findAuthenticatorDataSource(ctx, client, name, key)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(authenticator.GetId())
	_ = d.Set("key", authenticator.GetKey())
	_ = d.Set("name", authenticator.GetName())
	_ = d.Set("status", authenticator.GetStatus())
	_ = d.Set("type", authenticator.GetType())

	// Extract settings from AdditionalProperties
	if authenticator.AdditionalProperties != nil {
		if settings, ok := authenticator.AdditionalProperties["settings"]; ok && settings != nil {
			b, _ := json.Marshal(settings)
			_ = d.Set("settings", string(b))
		}

		// Extract provider from AdditionalProperties
		if providerRaw, ok := authenticator.AdditionalProperties["provider"]; ok && providerRaw != nil {
			providerMap, ok := providerRaw.(map[string]interface{})
			if ok {
				b, _ := json.Marshal(providerMap)
				dataMap := map[string]interface{}{}
				_ = json.Unmarshal([]byte(string(b)), &dataMap)
				b, _ = json.Marshal(dataMap)
				_ = d.Set("provider_json", string(b))

				if provType, ok := providerMap["type"].(string); ok {
					_ = d.Set("provider_type", provType)
				}

				if config, ok := providerMap["configuration"].(map[string]interface{}); ok {
					// Extract configuration based on authenticator type
					switch authenticator.GetType() {
					case "security_key":
						if hostname, ok := config["hostName"].(string); ok {
							_ = d.Set("provider_hostname", hostname)
						}
						if authPort, ok := config["authPort"].(float64); ok {
							_ = d.Set("provider_auth_port", int(authPort))
						}
						if instanceId, ok := config["instanceId"].(string); ok {
							_ = d.Set("provider_instance_id", instanceId)
						}
					default:
						logger(meta).Debug("Authenticator type using standard data source configuration",
							"authenticator_type", authenticator.GetType())
					}

					// Extract configuration based on provider type
					if provType, ok := providerMap["type"].(string); ok {
						switch provType {
						case "DUO":
							if host, ok := config["host"].(string); ok {
								_ = d.Set("provider_host", host)
							}
							if secretKey, ok := config["secretKey"].(string); ok {
								_ = d.Set("provider_secret_key", secretKey)
							}
							if integrationKey, ok := config["integrationKey"].(string); ok {
								_ = d.Set("provider_integration_key", integrationKey)
							}
						default:
							logger(meta).Debug("Unknown provider type in data source - using default behavior",
								"provider_type", provType,
								"authenticator_key", authenticator.GetKey(),
								"supported_types", []string{"DUO"})
						}
					}

					// Extract common template configuration
					if template, ok := config["userNameTemplate"].(map[string]interface{}); ok {
						if t, ok := template["template"].(string); ok {
							_ = d.Set("provider_user_name_template", t)
						}
					}
				}
			}
		}
	}
	return nil
}

func findAuthenticatorDataSource(ctx context.Context, client *v6okta.APIClient, name, key string) (*v6okta.AuthenticatorBase, error) {
	authenticators, _, err := client.AuthenticatorAPI.ListAuthenticators(ctx).Execute()
	if err != nil {
		return nil, err
	}
	for _, authUnion := range authenticators {
		// Convert union type to AuthenticatorBase via JSON marshaling
		authBytes, err := json.Marshal(authUnion)
		if err != nil {
			continue
		}
		var authenticator v6okta.AuthenticatorBase
		if err := json.Unmarshal(authBytes, &authenticator); err != nil {
			continue
		}

		if key != "custom_otp" {
			if authenticator.GetName() == name {
				return &authenticator, nil
			}
			if authenticator.GetKey() == key {
				return &authenticator, nil
			}
		} else {
			if authenticator.GetName() == name && authenticator.GetKey() == key {
				return &authenticator, nil
			} else {
				return nil, fmt.Errorf("authenticator with name '%s' and/or key '%s' does not exist", name, key)
			}
		}
	}
	if key != "" {
		return nil, fmt.Errorf("authenticator with key '%s' does not exist", key)
	}
	return nil, fmt.Errorf("authenticator with name '%s' does not exist", name)
}
