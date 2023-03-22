package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func dataSourceAuthServerClaim() *schema.Resource {
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
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"scopes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Auth server claim list of scopes",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"claim_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"always_include_in_token": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceAuthServerClaimRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		claim, _, err = getOktaClientFromMetadata(m).AuthorizationServer.GetOAuth2Claim(ctx, d.Get("auth_server_id").(string), id)
	} else {
		claim, err = getAuthServerClaimByName(ctx, m, d.Get("auth_server_id").(string), name)
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
		_ = d.Set("scopes", convertStringSliceToSet(claim.Conditions.Scopes))
	}
	return nil
}

func getAuthServerClaimByName(ctx context.Context, m interface{}, authServerID, name string) (*sdk.OAuth2Claim, error) {
	claims, _, err := getOktaClientFromMetadata(m).AuthorizationServer.ListOAuth2Claims(ctx, authServerID)
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
