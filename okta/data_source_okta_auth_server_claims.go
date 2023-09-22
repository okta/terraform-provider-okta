package okta

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func dataSourceAuthServerClaims() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthServerClaimsRead,
		Schema: map[string]*schema.Schema{
			"auth_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server ID",
			},
			"claims": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of authorization server claims retrieved from Okta with the following properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the claim.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the claim.",
						},
						"scopes": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Specifies the scopes for this Claim.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the claim.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value of the claim.",
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies whether the Claim is an Okta EL expression (`EXPRESSION`), a set of groups (`GROUPS`), or a system claim (`SYSTEM`)",
						},
						"claim_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies whether the Claim is for an access token (`RESOURCE`) or ID token (`IDENTITY`).",
						},
						"always_include_in_token": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies whether to include Claims in the token.",
						},
					},
				},
			},
		},
		Description: "Get a list of authorization server claims from Okta.",
	}
}

func dataSourceAuthServerClaimsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	claims, _, err := getOktaClientFromMetadata(m).AuthorizationServer.ListOAuth2Claims(ctx, d.Get("auth_server_id").(string))
	if err != nil {
		return diag.Errorf("failed to list authorization server claims: %v", err)
	}
	var s string
	arr := make([]map[string]interface{}, len(claims))
	for i := range claims {
		s += claims[i].Name
		arr[i] = flattenClaim(claims[i])
	}
	_ = d.Set("claims", arr)
	d.SetId(fmt.Sprintf("%s.%d", d.Get("auth_server_id").(string), crc32.ChecksumIEEE([]byte(s))))
	return nil
}

func flattenClaim(c *sdk.OAuth2Claim) map[string]interface{} {
	m := map[string]interface{}{
		"id":                      c.Id,
		"name":                    c.Name,
		"status":                  c.Status,
		"value":                   c.Value,
		"value_type":              c.ValueType,
		"claim_type":              c.ClaimType,
		"always_include_in_token": c.AlwaysIncludeInToken,
	}
	if c.Conditions != nil && len(c.Conditions.Scopes) > 0 {
		m["scopes"] = convertStringSliceToSet(c.Conditions.Scopes)
	}
	return m
}
