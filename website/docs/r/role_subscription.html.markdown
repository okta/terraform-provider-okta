---
layout: 'okta'
page_title: 'Okta: okta_role_subscription'
sidebar_current: 'docs-okta-resource-role-subscription'
description: |-
  Manages group subscription.
---

# okta_role_subscription

This resource allows you to configure subscriptions of a Role with a specific type.

## Example Usage

```hcl
resource "okta_role_subscription" "test" {
  role_type         = "SUPER_ADMIN"
  notification_type = "APP_IMPORT"
  status            = "unsubscribed"
}
```

## Argument Reference

- `role_type` - (Required) Type of the role. Valid values: `"SUPER_ADMIN"`, `"ORG_ADMIN"`, `"APP_ADMIN"`, `"USER_ADMIN"`,
  `"HELP_DESK_ADMIN"`, `"READ_ONLY_ADMIN"`, `"MOBILE_ADMIN"`, `"API_ADMIN"`, `"REPORT_ADMIN"`,
  `"GROUP_MEMBERSHIP_ADMIN"`.

- `notification_type` - (Required) Type of the notification. Valid values: `"CONNECTOR_AGENT"`, `"USER_LOCKED_OUT"`,
  `"APP_IMPORT"`, `"LDAP_AGENT"`, `"AD_AGENT"`, `"OKTA_ANNOUNCEMENT"`, `"OKTA_ISSUE"`, `"OKTA_UPDATE"`, `"IWA_AGENT"`,
  `"USER_DEPROVISION"`, `"REPORT_SUSPICIOUS_ACTIVITY"`, `"RATELIMIT_NOTIFICATION"`.

- `status` - (Optional) Subscription status. Valid values: `"subscribed"`, `"unsubscribed"`.

## Attributes Reference

- `id` - ID of the resource. Same a `notification_type`.

## Import

A role subscription can be imported via the Okta ID.

```
$ terraform import okta_role_subscription.example <role_type>/<notification_type>
```