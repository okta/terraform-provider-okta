package okta

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func dataSourceAuthenticator() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthenticatorRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key", "name"},
			},
			"key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "key"},
			},
			"settings": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Authenticator settings in JSON format",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the authenticator",
			},
			"provider_hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server host name or IP address",
			},
			"provider_auth_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The RADIUS server port (for example 1812). This is defined when the On-Prem RADIUS server is configured",
			},
			"provider_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_user_name_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Format expected by the provider",
			},
		},
	}
}

func dataSourceAuthenticatorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	name := d.Get("name").(string)
	key := d.Get("key").(string)
	if id == "" && name == "" && key == "" {
		return diag.Errorf("config must provide either 'id', 'name' or 'key' to retrieve the authenticator")
	}
	var (
		authenticator *okta.Authenticator
		err           error
	)
	if id != "" {
		authenticator, _, err = getOktaClientFromMetadata(m).Authenticator.GetAuthenticator(ctx, id)
	} else {
		authenticator, err = findAuthenticator(ctx, m, name, key)
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
		_ = d.Set("provider_type", authenticator.Provider.Type)
		_ = d.Set("provider_hostname", authenticator.Provider.Configuration.HostName)
		_ = d.Set("provider_auth_port", authenticator.Provider.Configuration.AuthPort)
		_ = d.Set("provider_instance_id", authenticator.Provider.Configuration.InstanceId)
		if authenticator.Provider.Configuration.UserNameTemplate != nil {
			_ = d.Set("provider_user_name_template", authenticator.Provider.Configuration.UserNameTemplate.Template)
		}
	}
	return nil
}

func findAuthenticator(ctx context.Context, m interface{}, name, key string) (*okta.Authenticator, error) {
	authenticators, _, err := getOktaClientFromMetadata(m).Authenticator.ListAuthenticators(ctx)
	if err != nil {
		return nil, err
	}
	for _, authenticator := range authenticators {
		if authenticator.Name == name {
			return authenticator, nil
		}
		if authenticator.Key == key {
			return authenticator, nil
		}
	}
	if key != "" {
		return nil, fmt.Errorf("authenticator with key '%s' does not exist", key)
	}
	return nil, fmt.Errorf("authenticator with name '%s' does not exist", name)
}
