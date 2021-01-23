package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAuthServerClaim() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthServerClaimCreate,
		ReadContext:   resourceAuthServerClaimRead,
		UpdateContext: resourceAuthServerClaimUpdate,
		DeleteContext: resourceAuthServerClaimDelete,
		Importer:      createNestedResourceImporter([]string{"auth_server_id", "id"}),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server claim name",
			},
			"auth_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server ID",
			},
			"scopes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Auth server claim list of scopes",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": statusSchema,
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{"EXPRESSION", "GROUPS", "SYSTEM"}),
				Default:          "EXPRESSION",
			},
			"claim_type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringInSlice([]string{"RESOURCE", "IDENTITY"}),
			},
			"always_include_in_token": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"group_filter_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{"STARTS_WITH", "EQUALS", "CONTAINS", "REGEX"}),
				Description:      "Required when value_type is GROUPS",
			},
		},
	}
}

func resourceAuthServerClaimCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	claim := buildAuthServerClaim(d)
	respClaim, _, err := getOktaClientFromMetadata(m).AuthorizationServer.CreateOAuth2Claim(ctx, d.Get("auth_server_id").(string), claim)
	if err != nil {
		return diag.Errorf("failed to create auth server claim: %v", err)
	}
	d.SetId(respClaim.Id)
	return resourceAuthServerClaimRead(ctx, d, m)
}

func resourceAuthServerClaimRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	claim, resp, err := getOktaClientFromMetadata(m).AuthorizationServer.GetOAuth2Claim(ctx, d.Get("auth_server_id").(string), d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get auth server: %v", err)
	}
	if claim == nil {
		d.SetId("")
		return nil
	}
	if claim.Conditions != nil && len(claim.Conditions.Scopes) > 0 {
		_ = d.Set("scopes", convertStringSetToInterface(claim.Conditions.Scopes))
	}
	_ = d.Set("name", claim.Name)
	_ = d.Set("status", claim.Status)
	_ = d.Set("value", claim.Value)
	_ = d.Set("value_type", claim.ValueType)
	_ = d.Set("claim_type", claim.ClaimType)
	_ = d.Set("always_include_in_token", claim.AlwaysIncludeInToken)
	_ = d.Set("group_filter_type", claim.GroupFilterType)
	return nil
}

func resourceAuthServerClaimUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	claim := buildAuthServerClaim(d)
	_, _, err := getOktaClientFromMetadata(m).AuthorizationServer.UpdateOAuth2Claim(ctx, d.Get("auth_server_id").(string), d.Id(), claim)
	if err != nil {
		return diag.Errorf("failed to update auth server claim: %v", err)
	}
	return resourceAuthServerClaimRead(ctx, d, m)
}

func resourceAuthServerClaimDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// 'valueType' of 'SYSTEM' implies a default claim and cannot be deleted.
	// System claims can be excluded from id tokens by changing the value of 'alwaysIncludeInToken'.
	if d.Get("value_type").(string) == "SYSTEM" && d.Get("always_include_in_token").(bool) {
		return nil
	}
	_, err := getOktaClientFromMetadata(m).AuthorizationServer.DeleteOAuth2Claim(ctx, d.Get("auth_server_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("failed to delete auth server claim: %v", err)
	}
	return nil
}

func buildAuthServerClaim(d *schema.ResourceData) okta.OAuth2Claim {
	return okta.OAuth2Claim{
		Status:               d.Get("status").(string),
		ClaimType:            d.Get("claim_type").(string),
		ValueType:            d.Get("value_type").(string),
		Value:                d.Get("value").(string),
		AlwaysIncludeInToken: boolPtr(d.Get("always_include_in_token").(bool)),
		Name:                 d.Get("name").(string),
		Conditions:           &okta.OAuth2ClaimConditions{Scopes: convertInterfaceToStringSetNullable(d.Get("scopes"))},
		GroupFilterType:      d.Get("group_filter_type").(string),
	}
}
