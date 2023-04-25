---
layout: 'okta'
page_title: 'Okta: okta_policy_rule_password'
sidebar_current: 'docs-okta-resource-policy-rule-password'
description: |-
  Creates a Password Policy Rule.
---

# okta_policy_rule_password

Creates a Password Policy Rule.

This resource allows you to create and configure a Password Policy Rule.

## Argument Reference

The following arguments are supported:
  
- `policy_id` - (Required) Policy ID.

- `name` - (Required) Policy Rule Name. Type `"string"`

- `priority` - (Optional) Policy Rule Priority, this attribute can be set to a valid priority. To avoid endless diff situation we error if an invalid priority is provided. API defaults it to the last (lowest) if not there. Type `"number"`

- `status` - (Optional) Policy Rule Status: `"ACTIVE"` or `"INACTIVE"`. Type `"string"`

- `password_change` - (Optional) Allow or deny a user to change their password: `"ALLOW"` or `"DENY"`. By default, it is `"ALLOW"`. Type `"string"`

- `password_reset` - (Optional) Allow or deny a user to reset their password: `"ALLOW"` or `"DENY"`. By default, it is `"ALLOW"`. Type `"string"`

- `password_unlock` - (Optional) Allow or deny a user to unlock: `"ALLOW"` or `"DENY"`. By default, it is `"DENY"`, Type `"string"`

- `network_connection` - (Optional) Network selection mode: `"ANYWHERE"`, `"ZONE"`, `"ON_NETWORK"`, or `"OFF_NETWORK"`. Type `"string"`

- `network_includes` - (Optional) The network zones to include. Conflicts with `network_excludes`. Type `"list(string)"`

- `network_excludes` - (Optional) The network zones to exclude. Conflicts with `network_includes`. Type `"list(string)"`

- `users_excluded` - (Optional) The users to exclude. Type `"set(string)"`

## Attributes Reference

- `id` - ID of the Rule.
  
- `policy_id` - Policy ID.

## Import

A Policy Rule can be imported via the Policy and Rule ID.

```
$ terraform import okta_policy_rule_password.example &#60;policy id&#62;/&#60;rule id&#62;
```
