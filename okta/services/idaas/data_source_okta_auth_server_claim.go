package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func DataSourceAuthServerClaim() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthServerClaimRead,
		Schema: map[string]*schema.Schema{
			"auth_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server ID",
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "Name of the claim. Conflicts with `name`.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
				Description:   "Name of the claim. Conflicts with `id`.",
			},
			"scopes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Auth server claim list of scopes",
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
		Description: "Get authorization server claim from Okta.",
	}
}

func dataSourceAuthServerClaimRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	name := d.Get("name").(string)
	if id == "" && name == "" {
		return diag.Errorf("config must provide either 'id' or 'name' to retrieve the auth server claim")
	}
	var (
		err   error
		claim *sdk.OAuth2Claim
	)
	if id != "" {
		claim, _, err = GetOktaClientFromMetadata(meta).AuthorizationServer.GetOAuth2Claim(ctx, d.Get("auth_server_id").(string), id)
	} else {
		claim, err = getAuthServerClaimByName(ctx, meta, d.Get("auth_server_id").(string), name)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(claim.Id)
	_ = d.Set("name", claim.Name)
	_ = d.Set("status", claim.Status)
	_ = d.Set("value", claim.Value)
	_ = d.Set("value_type", claim.ValueType)
	_ = d.Set("claim_type", claim.ClaimType)
	_ = d.Set("always_include_in_token", claim.AlwaysIncludeInToken)
	if claim.Conditions != nil && len(claim.Conditions.Scopes) > 0 {
		_ = d.Set("scopes", utils.ConvertStringSliceToSet(claim.Conditions.Scopes))
	}
	return nil
}

func getAuthServerClaimByName(ctx context.Context, meta interface{}, authServerID, name string) (*sdk.OAuth2Claim, error) {
	claims, _, err := GetOktaClientFromMetadata(meta).AuthorizationServer.ListOAuth2Claims(ctx, authServerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list authorization server claims: %v", err)
	}
	for i := range claims {
		if claims[i].Name == name {
			return claims[i], nil
		}
	}
	return nil, fmt.Errorf("auth server claim with name '%s' does not exist", name)
}
