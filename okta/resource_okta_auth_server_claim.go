package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAuthServerClaim() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthServerClaimCreate,
		ReadContext:   resourceAuthServerClaimRead,
		UpdateContext: resourceAuthServerClaimUpdate,
		DeleteContext: resourceAuthServerClaimDelete,
		Importer:      createNestedResourceImporter([]string{"auth_server_id", "id"}),
		Description:   "Creates an Authorization Server Claim. This resource allows you to create and configure an Authorization Server Claim.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the claim.",
			},
			"auth_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the authorization server.",
			},
			"scopes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The list of scopes the auth server claim is tied to.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": statusSchema,
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of the claim.",
			},
			"value_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "EXPRESSION",
				Description: "The type of value of the claim. It can be set to `EXPRESSION` or `GROUPS`. It defaults to `EXPRESSION`.",
			},
			"claim_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies whether the claim is for an access token `RESOURCE` or ID token `IDENTITY`.",
			},
			"always_include_in_token": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specifies whether to include claims in token, by default it is set to `true`.",
			},
			"group_filter_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type of group filter if `value_type` is `GROUPS`. Can be set to one of the following `STARTS_WITH`, `EQUALS`, `CONTAINS`, `REGEX`.",
			},
		},
	}
}

func resourceAuthServerClaimCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	claim := buildAuthServerClaim(d)
	respClaim, _, err := getOktaClientFromMetadata(meta).AuthorizationServer.CreateOAuth2Claim(ctx, d.Get("auth_server_id").(string), claim)
	if err != nil {
		return diag.Errorf("failed to create auth server claim: %v", err)
	}
	d.SetId(respClaim.Id)
	return resourceAuthServerClaimRead(ctx, d, meta)
}

func resourceAuthServerClaimRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	claim, resp, err := getOktaClientFromMetadata(meta).AuthorizationServer.GetOAuth2Claim(ctx, d.Get("auth_server_id").(string), d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get auth server claim: %v", err)
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
	_ = d.Set("group_filter_type", claim.GroupFilterType)
	return nil
}

func resourceAuthServerClaimUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	claim := buildAuthServerClaim(d)
	_, _, err := getOktaClientFromMetadata(meta).AuthorizationServer.UpdateOAuth2Claim(ctx, d.Get("auth_server_id").(string), d.Id(), claim)
	if err != nil {
		return diag.Errorf("failed to update auth server claim: %v", err)
	}
	return resourceAuthServerClaimRead(ctx, d, meta)
}

func resourceAuthServerClaimDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// 'valueType' of 'SYSTEM' implies a default claim and cannot be deleted.
	// System claims can be excluded from id tokens by changing the value of 'alwaysIncludeInToken'.
	if d.Get("value_type").(string) == "SYSTEM" && d.Get("always_include_in_token").(bool) {
		return nil
	}
	_, err := getOktaClientFromMetadata(meta).AuthorizationServer.DeleteOAuth2Claim(ctx, d.Get("auth_server_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("failed to delete auth server claim: %v", err)
	}
	return nil
}

func buildAuthServerClaim(d *schema.ResourceData) sdk.OAuth2Claim {
	return sdk.OAuth2Claim{
		Status:               d.Get("status").(string),
		ClaimType:            d.Get("claim_type").(string),
		ValueType:            d.Get("value_type").(string),
		Value:                d.Get("value").(string),
		AlwaysIncludeInToken: boolPtr(d.Get("always_include_in_token").(bool)),
		Name:                 d.Get("name").(string),
		Conditions:           &sdk.OAuth2ClaimConditions{Scopes: convertInterfaceToStringSetNullable(d.Get("scopes"))},
		GroupFilterType:      d.Get("group_filter_type").(string),
	}
}
