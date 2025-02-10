package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRoleSubscription() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRoleSubscriptionRead,
		Schema: map[string]*schema.Schema{
			"role_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the role",
			},
			"notification_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the notification",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of subscription",
			},
		},
		Description: "Get subscriptions of a Role with a specific type",
	}
}

func dataSourceRoleSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	subscription, _, err := getOktaClientFromMetadata(meta).Subscription.GetRoleSubscriptionByNotificationType(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
	if err != nil {
		return diag.Errorf("failed get subscription: %v", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", d.Get("role_type").(string), d.Get("notification_type").(string)))
	_ = d.Set("status", subscription.Status)
	return nil
}
