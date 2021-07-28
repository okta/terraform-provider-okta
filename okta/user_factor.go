package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func handleUserFactorLifecycle(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getSupplementFromMetadata(m)
	if d.Get("status").(string) == statusActive {
		_, err := client.ActivateUserFactor(ctx, d.Get("user_id").(string), d.Id())
		if err != nil {
			return diag.Errorf("failed to activate user factor: %v", err)
		}
		return nil
	}
	_, err := client.DeactivateUserFactor(ctx, d.Get("user_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("failed to deactivate user factor: %v", err)
	}
	return nil
}
