package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEventHookVerification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventHookVerificationCreate,
		ReadContext:   resourceEventHookVerificationRead,
		DeleteContext: resourceEventHookVerificationDelete,
		Importer:      nil,
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

func resourceEventHookVerificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getOktaClientFromMetadata(m).EventHook.VerifyEventHook(ctx, d.Get("event_hook_id").(string))
	if err != nil {
		return diag.Errorf("failed to verify event hook sender: %v", err)
	}
	d.SetId(d.Get("event_hook_id").(string))
	return nil
}

func resourceEventHookVerificationRead(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}

func resourceEventHookVerificationDelete(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}
