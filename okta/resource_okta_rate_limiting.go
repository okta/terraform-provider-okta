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
		Description: `Manages rate limiting.
This resource allows you to configure the client-based rate limit and rate limiting communications settings.
~> **WARNING:** This resource is available only when using a SSWS API token in the provider config, it is incompatible with OAuth 2.0 authentication.
~> **WARNING:** This resource makes use of an internal/private Okta API endpoint that could change without notice rendering this resource inoperable.`,
		Schema: map[string]*schema.Schema{
			"login": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Called when accessing the Okta hosted login page. Valid values: `ENFORCE` _(Enforce limit and log per client (recommended))_, `DISABLE` _(Do nothing (not recommended))_, `PREVIEW` _(Log per client)_.",
			},
			"authorize": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Called during authentication. Valid values: `ENFORCE` _(Enforce limit and log per client (recommended))_, `DISABLE` _(Do nothing (not recommended))_, `PREVIEW` _(Log per client)_.",
			},
			"communications_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable or disable rate limiting communications. By default, it is `true`.",
				Default:     true,
			},
		},
	}
}

func resourceRateLimitingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, _, err := getAPISupplementFromMetadata(meta).SetClientBasedRateLimiting(ctx, buildRateLimiter(d))
	if err != nil {
		return diag.Errorf("failed to set client-based rate limiting: %v", err)
	}
	_, _, err = getAPISupplementFromMetadata(meta).SetRateLimitingCommunications(ctx, buildRateLimitingCommunications(d))
	if err != nil {
		return diag.Errorf("failed to set rate limiting communications: %v", err)
	}
	d.SetId("rate_limiting")
	return nil
}

func resourceRateLimitingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rl, _, err := getAPISupplementFromMetadata(meta).GetClientBasedRateLimiting(ctx)
	if err != nil || rl.GranularModeSettings == nil {
		return diag.Errorf("failed to get client-based rate limiting: %v", err)
	}
	_ = d.Set("login", rl.GranularModeSettings.LoginPage)
	_ = d.Set("authorize", rl.GranularModeSettings.OAuth2Authorize)
	comm, _, err := getAPISupplementFromMetadata(meta).GetRateLimitingCommunications(ctx)
	if err != nil {
		return diag.Errorf("failed to get rate limiting communications: %v", err)
	}
	_ = d.Set("communications_enabled", comm.RateLimitNotification)
	d.SetId("rate_limiting")
	return nil
}

func resourceRateLimitingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, _, err := getAPISupplementFromMetadata(meta).SetClientBasedRateLimiting(ctx, buildRateLimiter(d))
	if err != nil {
		return diag.Errorf("failed to set client-based rate limiting: %v", err)
	}
	_, _, err = getAPISupplementFromMetadata(meta).SetRateLimitingCommunications(ctx, buildRateLimitingCommunications(d))
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
