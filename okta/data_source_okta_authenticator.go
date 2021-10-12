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
				Type:     schema.TypeString,
				Optional: true,
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Authenticator settings in JSON format",
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        normalizeDataJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the authenticator. When specified in the terraform resource, will act as a filter when searching for the authenticator",
			},
		},
	}
}

func dataSourceAuthenticatorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return findDataSourceAuthenticator(ctx, d.Get("type").(string), d, m)
}

func findDataSourceAuthenticator(ctx context.Context, _type string, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authenticator, err := findOktaAuthenticator(ctx, _type, m)

	if err != nil {
		return diag.Errorf("failed to list authenticators: %+v", err)
	}

	d.SetId(authenticator.Id)
	_ = d.Set("key", authenticator.Key)
	_ = d.Set("name", authenticator.Name)
	b, err := json.Marshal(authenticator.Settings)
	if err == nil {
		_ = d.Set("settings", string(b))
	}
	_ = d.Set("status", authenticator.Status)
	_ = d.Set("type", authenticator.Type)

	return nil

}

func findOktaAuthenticator(ctx context.Context, _type string, m interface{}) (*okta.Authenticator, error) {
	// NOTE when okta-sdk-golang supports getting authenticator by ID search by
	// ID or search by Type such as is done in data_source_okta_group.go

	authenticators, _, err := getOktaClientFromMetadata(m).Authenticator.ListAuthenticators(ctx)
	if err != nil {
		return nil, err
	}
	for _, authenticator := range authenticators {
		if authenticator.Type == _type {
			return authenticator, nil
		}
	}
	return nil, fmt.Errorf("authenticator with type %q does not exist", _type)
}
