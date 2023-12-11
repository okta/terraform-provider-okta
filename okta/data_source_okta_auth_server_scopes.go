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
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of authorization server scopes retrieved from Okta with the following properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the Scope",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Scope",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the Scope",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the end user displayed in a consent dialog box",
						},
						"consent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates whether a consent dialog is needed for the Scope",
						},
						"metadata_publish": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the Scope should be included in the metadata",
						},
						"default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the Scope is a default Scope",
						},
						"system": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether Okta created the Scope",
						},
						"optional": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the Scope is optional",
						},
					},
				},
			},
		},
		Description: "Get a list of authorization server scopes from Okta.",
	}
}

func dataSourceAuthServerScopesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scopes, _, err := getOktaV3ClientFromMetadata(m).AuthorizationServerAPI.ListOAuth2Scopes(ctx, d.Get("auth_server_id").(string)).Execute()
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
