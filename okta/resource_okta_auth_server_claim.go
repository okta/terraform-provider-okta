package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
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
				ValidateDiagFunc: stringInSlice([]string{"EXPRESSION", "GROUPS"}),
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
	authServerClaim := buildAuthServerClaim(d)
	responseAuthServerClaim, _, err := getSupplementFromMetadata(m).CreateAuthorizationServerClaim(ctx, d.Get("auth_server_id").(string), *authServerClaim, nil)
	if err != nil {
		return diag.Errorf("failed to create auth server claim: %v", err)
	}
	d.SetId(responseAuthServerClaim.ID)
	return resourceAuthServerClaimRead(ctx, d, m)
}

func resourceAuthServerClaimRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authServerClaim, resp, err := getSupplementFromMetadata(m).GetAuthorizationServerClaim(ctx, d.Get("auth_server_id").(string), d.Id(), sdk.AuthorizationServerClaim{})
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get auth server: %v", err)
	}
	if authServerClaim == nil {
		d.SetId("")
		return nil
	}
	if authServerClaim.Conditions != nil && len(authServerClaim.Conditions.Scopes) > 0 {
		_ = d.Set("scopes", convertStringSetToInterface(authServerClaim.Conditions.Scopes))
	}
	_ = d.Set("name", authServerClaim.Name)
	_ = d.Set("status", authServerClaim.Status)
	_ = d.Set("value", authServerClaim.Value)
	_ = d.Set("value_type", authServerClaim.ValueType)
	_ = d.Set("claim_type", authServerClaim.ClaimType)
	_ = d.Set("always_include_in_token", authServerClaim.AlwaysIncludeInToken)
	_ = d.Set("group_filter_type", authServerClaim.GroupFilterType)
	return nil
}

func resourceAuthServerClaimUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authServerClaim := buildAuthServerClaim(d)
	_, _, err := getSupplementFromMetadata(m).UpdateAuthorizationServerClaim(ctx, d.Get("auth_server_id").(string), d.Id(), *authServerClaim, nil)
	if err != nil {
		return diag.Errorf("failed to update auth server claim: %v", err)
	}
	return resourceAuthServerClaimRead(ctx, d, m)
}

func resourceAuthServerClaimDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := getSupplementFromMetadata(m).DeleteAuthorizationServerClaim(ctx, d.Get("auth_server_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("failed to delete auth server claim: %v", err)
	}
	return nil
}

func buildAuthServerClaim(d *schema.ResourceData) *sdk.AuthorizationServerClaim {
	return &sdk.AuthorizationServerClaim{
		Status:               d.Get("status").(string),
		ClaimType:            d.Get("claim_type").(string),
		ValueType:            d.Get("value_type").(string),
		Value:                d.Get("value").(string),
		AlwaysIncludeInToken: d.Get("always_include_in_token").(bool),
		Name:                 d.Get("name").(string),
		Conditions:           &sdk.ClaimConditions{Scopes: convertInterfaceToStringSetNullable(d.Get("scopes"))},
		GroupFilterType:      d.Get("group_filter_type").(string),
	}
}
