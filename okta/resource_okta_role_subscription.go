package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRoleSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleSubscriptionCreate,
		ReadContext:   resourceRoleSubscriptionRead,
		UpdateContext: resourceRoleSubscriptionUpdate,
		DeleteContext: resourceRoleSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid role subscription specifier, expecting {roleType}/{notificationType}")
				}
				_ = d.Set("role_type", parts[0])
				_ = d.Set("notification_type", parts[1])
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"role_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: elemInSlice(append(validAdminRoles, "API_ADMIN")),
				Description:      "Type of the role",
			},
			"notification_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: elemInSlice(validNotificationTypes),
				Description:      "Type of the notification",
			},
			"status": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: elemInSlice([]string{"subscribed", "unsubscribed"}),
				Description:      "Status of subscription",
			},
		},
	}
}

func resourceRoleSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	status, ok := d.GetOk("status")
	if !ok {
		return resourceRoleSubscriptionRead(ctx, d, m)
	}
	subscription, _, err := getSupplementFromMetadata(m).GetRoleTypeSubscription(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
	if err != nil {
		return diag.Errorf("failed get subscription: %v", err)
	}
	if subscription.Status != status.(string) {
		if status == "subscribed" {
			_, err = getSupplementFromMetadata(m).RoleTypeSubscribe(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
		} else {
			_, err = getSupplementFromMetadata(m).RoleTypeUnsubscribe(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
		}
		if err != nil {
			return diag.Errorf("failed to change subscription: %v", err)
		}
	}
	d.SetId(d.Get("notification_type").(string))
	return resourceRoleSubscriptionRead(ctx, d, m)
}

func resourceRoleSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	subscription, _, err := getSupplementFromMetadata(m).GetRoleTypeSubscription(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
	if err != nil {
		return diag.Errorf("failed get subscription: %v", err)
	}
	_ = d.Set("status", subscription.Status)
	return nil
}

func resourceRoleSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus == newStatus {
		return nil
	}
	if newStatus == "subscribed" {
		_, err = getSupplementFromMetadata(m).RoleTypeSubscribe(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
	} else {
		_, err = getSupplementFromMetadata(m).RoleTypeUnsubscribe(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
	}
	if err != nil {
		return diag.Errorf("failed to change subscription: %v", err)
	}
	return nil
}

func resourceRoleSubscriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
