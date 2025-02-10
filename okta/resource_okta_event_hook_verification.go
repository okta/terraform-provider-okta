package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEventHookVerification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventHookVerificationCreate,
		ReadContext:   resourceFuncNoOp,
		DeleteContext: resourceFuncNoOp,
		Importer:      nil,
		Description:   "Verifies the Event Hook. The resource won't be created unless the URI provided in the event hook returns a valid JSON object with verification. See [Event Hooks](https://developer.okta.com/docs/concepts/event-hooks/#one-time-verification-request) documentation for details.",
		Schema: map[string]*schema.Schema{
			"event_hook_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Event hook ID",
			},
		},
	}
}

func resourceEventHookVerificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, _, err := getOktaClientFromMetadata(meta).EventHook.VerifyEventHook(ctx, d.Get("event_hook_id").(string))
	if err != nil {
		return diag.Errorf("failed to verify event hook sender: %v", err)
	}
	d.SetId(d.Get("event_hook_id").(string))
	return nil
}
