package idaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func resourceRoleSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleSubscriptionCreate,
		ReadContext:   resourceRoleSubscriptionRead,
		UpdateContext: resourceRoleSubscriptionUpdate,
		DeleteContext: utils.ResourceFuncNoOp,
		Description: `Manages group subscription.
		
This resource allows you to configure subscriptions of a Role with a specific type. 
Check [configure email notifications](https://help.okta.com/oie/en-us/Content/Topics/Security/custom-admin-role/administrator-email-settings.htm) 
page regarding what notifications are available for specific admin roles.`,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// https://developer.okta.com/docs/reference/api/admin-notifications/#role-types
				Description: `Type of the role. Valid values:
	'API_ADMIN',
	'APP_ADMIN',
	'CUSTOM',
	'GROUP_MEMBERSHIP_ADMIN',
	'HELP_DESK_ADMIN',
	'MOBILE_ADMIN',
	'ORG_ADMIN',
	'READ_ONLY_ADMIN',
	'REPORT_ADMIN',
	'SUPER_ADMIN',
	'USER_ADMIN'
	. See [API docs](https://developer.okta.com/docs/reference/api/admin-notifications/#role-types).`,
			},
			"notification_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `Type of the notification. Valid values: 
	- 'CONNECTOR_AGENT' -  Disconnects and reconnects: On-prem provisioning, on-prem MFA agents, and RADIUS server agent.
	- 'USER_LOCKED_OUT' - User lockouts.
	- 'APP_IMPORT' - App user import status.
	- 'LDAP_AGENT' - Disconnects and reconnects: LDAP agent.
	- 'AD_AGENT' - Disconnects and reconnects: AD agent.
	- 'OKTA_ANNOUNCEMENT' - Okta release notes and announcements.
	- 'OKTA_UPDATE' - Scheduled system updates.
	- 'IWA_AGENT' - Disconnects and reconnects: IWA agent.
	- 'USER_DEPROVISION' - User deprovisions.
	- 'REPORT_SUSPICIOUS_ACTIVITY' - User reporting of suspicious activity.
	- 'RATELIMIT_NOTIFICATION' - Rate limit warning and violation.
	- 'AGENT_AUTO_UPDATE_NOTIFICATION' - Agent auto-update notifications: AD Agent.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subscription status. Valid values: `subscribed`, `unsubscribed`.",
			},
		},
	}
}

func resourceRoleSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validateSubscriptions(d.Get("role_type").(string), d.Get("notification_type").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	status, ok := d.GetOk("status")
	if !ok {
		return resourceRoleSubscriptionRead(ctx, d, meta)
	}
	subscription, _, err := getOktaClientFromMetadata(meta).Subscription.GetRoleSubscriptionByNotificationType(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
	if err != nil {
		return diag.Errorf("failed get subscription: %v", err)
	}
	if subscription.Status != status.(string) {
		if status == "subscribed" {
			_, err = getOktaClientFromMetadata(meta).Subscription.SubscribeRoleSubscriptionByNotificationType(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
		} else {
			_, err = getOktaClientFromMetadata(meta).Subscription.UnsubscribeRoleSubscriptionByNotificationType(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
		}
		if err != nil {
			return diag.Errorf("failed to change subscription: %v", err)
		}
	}
	d.SetId(d.Get("notification_type").(string))
	return resourceRoleSubscriptionRead(ctx, d, meta)
}

func resourceRoleSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	subscription, _, err := getOktaClientFromMetadata(meta).Subscription.GetRoleSubscriptionByNotificationType(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
	if err != nil {
		return diag.Errorf("failed get subscription: %v", err)
	}
	if subscription == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("status", subscription.Status)
	return nil
}

func resourceRoleSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validateSubscriptions(d.Get("role_type").(string), d.Get("notification_type").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus == newStatus {
		return nil
	}
	if newStatus == "subscribed" {
		_, err = getOktaClientFromMetadata(meta).Subscription.SubscribeRoleSubscriptionByNotificationType(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
	} else {
		_, err = getOktaClientFromMetadata(meta).Subscription.UnsubscribeRoleSubscriptionByNotificationType(ctx, d.Get("role_type").(string), d.Get("notification_type").(string))
	}
	if err != nil {
		return diag.Errorf("failed to change subscription: %v", err)
	}
	return nil
}

func validateSubscriptions(role, notification string) error {
	switch {
	case notification == "CONNECTOR_AGENT" || notification == "APP_IMPORT" || notification == "LDAP_AGENT" ||
		notification == "AD_AGENT" || notification == "IWA_AGENT":
		if role == "SUPER_ADMIN" || role == "ORG_ADMIN" || role == "APP_ADMIN" {
			return nil
		}
	case notification == "USER_LOCKED_OUT":
		if role == "SUPER_ADMIN" || role == "ORG_ADMIN" || role == "USER_ADMIN" || role == "HELP_DESK_ADMIN" {
			return nil
		}
	case notification == "USER_DEPROVISION":
		if role == "SUPER_ADMIN" || role == "MOBILE_ADMIN" || role == "APP_ADMIN" || role == "API_ACCESS_MANAGEMENT_ADMIN" {
			return nil
		}
	case notification == "REPORT_SUSPICIOUS_ACTIVITY":
		if role == "SUPER_ADMIN" || role == "ORG_ADMIN" {
			return nil
		}
	case notification == "RATELIMIT_NOTIFICATION" || notification == "AGENT_AUTO_UPDATE_NOTIFICATION":
		if role == "SUPER_ADMIN" {
			return nil
		}
	case notification == "OKTA_ANNOUNCEMENT" || notification == "OKTA_ISSUE" || notification == "OKTA_UPDATE":
		return nil
	}
	return fmt.Errorf("'%s' notification is not aplicable for the '%s' role", notification, role)
}
