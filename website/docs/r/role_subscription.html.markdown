---
layout: 'okta'
page_title: 'Okta: okta_role_subscription'
sidebar_current: 'docs-okta-resource-role-subscription'
description: |-
  Manages group subscription.
---

# okta_role_subscription

This resource allows you to configure subscriptions of a Role with a specific type. 
Check [configure email notifications](https://help.okta.com/oie/en-us/Content/Topics/Security/custom-admin-role/administrator-email-settings.htm) 
page regarding what notifications are available for specific admin roles.

## Example Usage

```hcl
resource "okta_role_subscription" "test" {
  role_type         = "SUPER_ADMIN"
  notification_type = "APP_IMPORT"
  status            = "unsubscribed"
}
```

## Argument Reference

- `role_type` - (Required) Type of the role. Valid values: `"API_ADMIN"`, `"APP_ADMIN"`, `"GROUP_MEMBERSHIP_ADMIN"`, `"HELP_DESK_ADMIN"`, `"MOBILE_ADMIN"`, `"ORG_ADMIN"`, `"READ_ONLY_ADMIN"`, `"REPORT_ADMIN"`, `"SUPER_ADMIN"`, `"USER_ADMIN"`.  See [API docs](https://developer.okta.com/docs/reference/api/admin-notifications/#role-types)

- `notification_type` - (Required) Type of the notification. Valid values: 
  - `"CONNECTOR_AGENT"` -  Disconnects and reconnects: On-prem provisioning, on-prem MFA agents, and RADIUS server agent.
  - `"USER_LOCKED_OUT"` - User lockouts.
  - `"APP_IMPORT"` - App user import status.
  - `"LDAP_AGENT"` - Disconnects and reconnects: LDAP agent.
  - `"AD_AGENT"` - Disconnects and reconnects: AD agent.
  - `"OKTA_ANNOUNCEMENT"` - Okta release notes and announcements.
  - `"OKTA_ISSUE"` - Trust incidents and updates.
  - `"OKTA_UPDATE"` - Scheduled system updates.
  - `"IWA_AGENT"` - Disconnects and reconnects: IWA agent.
  - `"USER_DEPROVISION"` - User deprovisions.
  - `"REPORT_SUSPICIOUS_ACTIVITY"` - User reporting of suspicious activity.
  - `"RATELIMIT_NOTIFICATION"` - Rate limit warning and violation.
  - `"AGENT_AUTO_UPDATE_NOTIFICATION"` - Agent auto-update notifications: AD Agent.

- `status` - (Optional) Subscription status. Valid values: `"subscribed"`, `"unsubscribed"`.

## Attributes Reference

- `id` - ID of the resource. Same a `notification_type`.

## Import

A role subscription can be imported via the Okta ID.

```
$ terraform import okta_role_subscription.example &#60;role_type&#62;/&#60;notification_type&#62;
```
