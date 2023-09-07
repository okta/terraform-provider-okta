package okta

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v3/okta"
)

func dataSourceAuthServerScopes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthServerScopesRead,
		Schema: map[string]*schema.Schema{
			"auth_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server ID",
			},
			"scopes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata_publish": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"system": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"optional": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAuthServerScopesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scopes, _, err := getOktaV3ClientFromMetadata(m).AuthorizationServerApi.ListOAuth2Scopes(ctx, d.Get("auth_server_id").(string)).Execute()
	if err != nil {
		return diag.Errorf("failed to list auth server scopes: %v", err)
	}
	var s string
	arr := make([]map[string]interface{}, len(scopes))
	for i := range scopes {
		s += scopes[i].GetName()
		arr[i] = flattenScope(scopes[i])
	}
	_ = d.Set("scopes", arr)
	d.SetId(fmt.Sprintf("%s.%d", d.Get("auth_server_id").(string), crc32.ChecksumIEEE([]byte(s))))
	return nil
}

func flattenScope(s okta.OAuth2Scope) map[string]interface{} {
	return map[string]interface{}{
		"id":               s.GetId(),
		"name":             s.GetName(),
		"description":      s.GetDescription(),
		"display_name":     s.GetDisplayName(),
		"consent":          s.GetConsent(),
		"metadata_publish": s.GetMetadataPublish(),
		"default":          s.GetDefault(),
		"system":           s.GetSystem(),
		"optional":         s.GetOptional(),
	}
}
