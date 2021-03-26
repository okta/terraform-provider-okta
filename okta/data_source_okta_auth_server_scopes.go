package okta

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
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
					},
				},
			},
		},
	}
}

func dataSourceAuthServerScopesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scopes, _, err := getOktaClientFromMetadata(m).AuthorizationServer.ListOAuth2Scopes(ctx, d.Get("auth_server_id").(string), nil)
	if err != nil {
		return diag.Errorf("failed to list auth server scopes: %v", err)
	}
	var s string
	arr := make([]map[string]interface{}, len(scopes))
	for i := range scopes {
		s += scopes[i].Name
		arr[i] = flattenScope(scopes[i])
	}
	_ = d.Set("scopes", arr)
	d.SetId(fmt.Sprintf("%s.%d", d.Get("auth_server_id").(string), crc32.ChecksumIEEE([]byte(s))))
	return nil
}

func flattenScope(s *okta.OAuth2Scope) map[string]interface{} {
	return map[string]interface{}{
		"id":               s.Id,
		"name":             s.Name,
		"description":      s.Description,
		"consent":          s.Consent,
		"metadata_publish": s.MetadataPublish,
		"default":          s.Default,
		"system":           s.System,
	}
}
