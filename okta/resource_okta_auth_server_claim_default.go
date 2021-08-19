package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAuthServerClaimDefault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthServerClaimDefaultUpdate,
		ReadContext:   resourceAuthServerClaimDefaultRead,
		UpdateContext: resourceAuthServerClaimDefaultUpdate,
		DeleteContext: resourceAuthServerClaimDefaultDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid resource import specifier, expecting the following format: <auth_server_id>/<id> or <auth_server_id>/<name>")
				}
				_ = d.Set("auth_server_id", parts[0])
				if contains(validDefaultAuthServerClaims, parts[1]) {
					c, err := findClaim(ctx, meta, parts[0], parts[1])
					if err != nil {
						return nil, err
					}
					d.SetId(c.Id)
				} else {
					d.SetId(parts[1])
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Default auth server claim name",
				ValidateDiagFunc: elemInSlice(validDefaultAuthServerClaims),
				ForceNew:         true,
			},
			"auth_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server ID",
				ForceNew:    true,
			},
			"scopes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Auth server claim list of scopes",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("name") != "sub"
				},
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

func resourceAuthServerClaimDefaultRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	claim, resp, err := getOktaClientFromMetadata(m).AuthorizationServer.GetOAuth2Claim(ctx, d.Get("auth_server_id").(string), d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get auth server default claim: %v", err)
	}
	if claim == nil {
		d.SetId("")
		return nil
	}
	if claim.Conditions != nil && len(claim.Conditions.Scopes) > 0 {
		_ = d.Set("scopes", convertStringSliceToSet(claim.Conditions.Scopes))
	}
	_ = d.Set("name", claim.Name)
	_ = d.Set("status", claim.Status)
	_ = d.Set("value", claim.Value)
	_ = d.Set("value_type", claim.ValueType)
	_ = d.Set("claim_type", claim.ClaimType)
	_ = d.Set("always_include_in_token", claim.AlwaysIncludeInToken)
	return nil
}

func resourceAuthServerClaimDefaultUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Id() == "" {
		claim, err := findClaim(ctx, m, d.Get("auth_server_id").(string), d.Get("name").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(claim.Id)
		if claim.Conditions != nil && len(claim.Conditions.Scopes) > 0 {
			_ = d.Set("scopes", convertStringSliceToSet(claim.Conditions.Scopes))
		}
		_ = d.Set("name", claim.Name)
		_ = d.Set("status", claim.Status)
		_ = d.Set("value_type", claim.ValueType)
		_ = d.Set("claim_type", claim.ClaimType)
		_ = d.Set("always_include_in_token", claim.AlwaysIncludeInToken)
		if d.Get("name").(string) != "sub" {
			_ = d.Set("value", claim.Value)
			return nil // all the values are computed, so just stop here
		}
	}
	if d.Get("name").(string) != "sub" {
		// all the default claims except "sub" are immutable
		return resourceAuthServerClaimDefaultRead(ctx, d, m)
	} else if d.Get("value").(string) == "" {
		return diag.Errorf("'value' is required parameter for 'sub' claim")
	}
	claim := buildAuthServerClaimDefault(d)
	_, _, err := getOktaClientFromMetadata(m).AuthorizationServer.UpdateOAuth2Claim(ctx, d.Get("auth_server_id").(string), d.Id(), claim)
	if err != nil {
		return diag.Errorf("failed to update auth server default claim: %v", err)
	}
	return resourceAuthServerClaimDefaultRead(ctx, d, m)
}

// Default claims are immutable.
func resourceAuthServerClaimDefaultDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func buildAuthServerClaimDefault(d *schema.ResourceData) okta.OAuth2Claim {
	return okta.OAuth2Claim{
		Status:               d.Get("status").(string),
		ClaimType:            d.Get("claim_type").(string),
		ValueType:            d.Get("value_type").(string),
		Value:                d.Get("value").(string),
		AlwaysIncludeInToken: boolPtr(d.Get("always_include_in_token").(bool)),
		Name:                 d.Get("name").(string),
		Conditions:           &okta.OAuth2ClaimConditions{Scopes: convertInterfaceToStringSetNullable(d.Get("scopes"))},
	}
}

func findClaim(ctx context.Context, m interface{}, serverID, name string) (*okta.OAuth2Claim, error) {
	claims, resp, err := getOktaClientFromMetadata(m).AuthorizationServer.ListOAuth2Claims(ctx, serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to list auth server claims: %v", err)
	}
	for {
		for _, claim := range claims {
			if claim.Name == name {
				return claim, nil
			}
		}
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &claims)
			if err != nil {
				return nil, fmt.Errorf("failed to auth server claims: %v", err)
			}
			continue
		} else {
			break
		}
	}
	return nil, fmt.Errorf("no claim '%s' found for auth server '%s'", name, serverID)
}

var validDefaultAuthServerClaims = []string{
	"sub", "address", "birthdate", "email", "email_verified",
	"family_name", "gender", "given_name", "locale", "middle_name", "name", "nickname", "phone_number",
	"picture", "preferred_username", "profile", "updated_at", "website", "zoneinfo",
}
