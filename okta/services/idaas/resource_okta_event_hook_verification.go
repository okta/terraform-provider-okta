package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func resourceEventHookVerification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventHookVerificationCreate,
		ReadContext:   resourceEventHookVerificationRead,
		UpdateContext: resourceEventHookVerificationUpdate,
		DeleteContext: utils.ResourceFuncNoOp,
		Importer:      nil,
		Description:   "Verifies the Event Hook. The resource won't be created unless the URI provided in the event hook returns a valid JSON object with verification. See [Event Hooks](https://developer.okta.com/docs/concepts/event-hooks/#one-time-verification-request) documentation for details.",
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			// When the API reports the hook as UNVERIFIED, force a planned diff so that
			// UpdateContext is invoked on the next apply to re-trigger verification.
			if d.Id() != "" && d.Get("verification_status").(string) == "UNVERIFIED" {
				return d.SetNew("verification_status", "VERIFIED")
			}
			return nil
		},
		Schema: map[string]*schema.Schema{
			"event_hook_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Event hook ID",
			},
			"verification_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Verification status of the Event hook",
			},
		},
	}
}

func resourceEventHookVerificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hook, _, err := getOktaV6ClientFromMetadata(meta).EventHookAPI.VerifyEventHook(ctx, d.Get("event_hook_id").(string)).Execute()
	if err != nil {
		return diag.Errorf("failed to verify event hook sender: %v", err)
	}
	d.SetId(d.Get("event_hook_id").(string))
	_ = d.Set("verification_status", hook.VerificationStatus)
	return nil
}

func resourceEventHookVerificationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hook, resp, err := getOktaV6ClientFromMetadata(meta).EventHookAPI.GetEventHook(ctx, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V6(resp, err); err != nil {
		return diag.Errorf("failed to get event hook: %v", err)
	}
	if hook == nil {
		d.SetId("")
		return nil
	}
	d.SetId(d.Get("event_hook_id").(string))
	_ = d.Set("verification_status", hook.VerificationStatus)
	return nil
}

func resourceEventHookVerificationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hook, _, err := getOktaV6ClientFromMetadata(meta).EventHookAPI.VerifyEventHook(ctx, d.Get("event_hook_id").(string)).Execute()
	if err != nil {
		return diag.Errorf("failed to verify event hook sender: %v", err)
	}
	d.SetId(d.Get("event_hook_id").(string))
	_ = d.Set("verification_status", hook.VerificationStatus)
	return nil
}
