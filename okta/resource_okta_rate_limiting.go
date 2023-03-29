package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceRateLimiting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRateLimitingCreate,
		ReadContext:   resourceRateLimitingRead,
		UpdateContext: resourceRateLimitingUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"login": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Called when accessing the Okta hosted login page.",
			},
			"authorize": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Called during authentication.",
			},
			"communications_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables rate limit warning, violation, notification emails and banners when this org meets rate limits.",
				Default:     true,
			},
		},
	}
}

func resourceRateLimitingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getSupplementFromMetadata(m).SetClientBasedRateLimiting(ctx, buildRateLimiter(d))
	if err != nil {
		return diag.Errorf("failed to set client-based rate limiting: %v", err)
	}
	_, _, err = getSupplementFromMetadata(m).SetRateLimitingCommunications(ctx, buildRateLimitingCommunications(d))
	if err != nil {
		return diag.Errorf("failed to set rate limiting communications: %v", err)
	}
	d.SetId("rate_limiting")
	return nil
}

func resourceRateLimitingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rl, _, err := getSupplementFromMetadata(m).GetClientBasedRateLimiting(ctx)
	if err != nil || rl.GranularModeSettings == nil {
		return diag.Errorf("failed to get client-based rate limiting: %v", err)
	}
	_ = d.Set("login", rl.GranularModeSettings.LoginPage)
	_ = d.Set("authorize", rl.GranularModeSettings.OAuth2Authorize)
	comm, _, err := getSupplementFromMetadata(m).GetRateLimitingCommunications(ctx)
	if err != nil {
		return diag.Errorf("failed to get rate limiting communications: %v", err)
	}
	_ = d.Set("communications_enabled", *comm.RateLimitNotification)
	d.SetId("rate_limiting")
	return nil
}

func resourceRateLimitingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getSupplementFromMetadata(m).SetClientBasedRateLimiting(ctx, buildRateLimiter(d))
	if err != nil {
		return diag.Errorf("failed to set client-based rate limiting: %v", err)
	}
	_, _, err = getSupplementFromMetadata(m).SetRateLimitingCommunications(ctx, buildRateLimitingCommunications(d))
	if err != nil {
		return diag.Errorf("failed to set rate limiting communications: %v", err)
	}
	return nil
}

func buildRateLimiter(d *schema.ResourceData) sdk.ClientRateLimitMode {
	return sdk.ClientRateLimitMode{
		Mode: "PREVIEW",
		GranularModeSettings: &sdk.RateLimitGranularModeSettings{
			OAuth2Authorize: d.Get("authorize").(string),
			LoginPage:       d.Get("login").(string),
		},
	}
}

func buildRateLimitingCommunications(d *schema.ResourceData) sdk.RateLimitingCommunications {
	return sdk.RateLimitingCommunications{
		RateLimitNotification: boolPtr(d.Get("communications_enabled").(bool)),
	}
}
