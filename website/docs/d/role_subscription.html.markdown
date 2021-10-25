---
layout: "okta"
page_title: "Okta: okta_role_subscription"
sidebar_current: "docs-okta-datasource-role-subscription"
description: |-
  Get subscriptions of a Role with a specific type
---

# okta_role_subscription

Use this data source to retrieve role subscription with a specific type.

## Example Usage

```hcl
data "okta_role_subscription" "example" {
  notification_type = "APP_IMPORT"
  role_type         = "SUPER_ADMIN"
}
```

## Arguments Reference

- `role_type` - (Required) Type of the role. Valid values: `"SUPER_ADMIN"`, `"ORG_ADMIN"`, `"APP_ADMIN"`, `"USER_ADMIN"`, 
  `"HELP_DESK_ADMIN"`, `"READ_ONLY_ADMIN"`, `"MOBILE_ADMIN"`, `"API_ACCESS_MANAGEMENT_ADMIN"`, `"REPORT_ADMIN"`, 
  `"GROUP_MEMBERSHIP_ADMIN"`.

- `notification_type` - (Required) Type of the notification. Valid values: `"CONNECTOR_AGENT"`, `"USER_LOCKED_OUT"`, 
  `"APP_IMPORT"`, `"LDAP_AGENT"`, `"AD_AGENT"`, `"OKTA_ANNOUNCEMENT"`, `"OKTA_ISSUE"`, `"OKTA_UPDATE"`, `"IWA_AGENT"`, 
  `"USER_DEPROVISION"`, `"REPORT_SUSPICIOUS_ACTIVITY"`, `"RATELIMIT_NOTIFICATION"`.

## Attributes Reference

- `id` - ID of the resource. Same a `notification_type`.

- `status` - Subscription status.
