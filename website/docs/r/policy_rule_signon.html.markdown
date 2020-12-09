---
layout: 'okta'
page_title: 'Okta: okta_policy_rule_signon'
sidebar_current: 'docs-okta-resource-policy-rule-signon'
description: |-
  Creates a Sign On Policy Rule.
---

# okta_policy_rule_signon

Creates a Sign On Policy Rule.

## Argument Reference

The following arguments are supported:

- `policyid` - (Required) Policy ID.

- `name` - (Required) Policy Rule Name.

- `priority` - (Optional) Policy Rule Priority, this attribute can be set to a valid priority. To avoid endless diff situation we error if an invalid priority is provided. API defaults it to the last/lowest if not there.

- `status` - (Optional) Policy Rule Status: `"ACTIVE"` or `"INACTIVE"`.

- `authtype` - (Optional) Authentication entrypoint: `"ANY"` or `"RADIUS"`.

- `access` - (Optional) Allow or deny access based on the rule conditions: `"ALLOW"` or `"DENY"`. The default is `"ALLOW"`.

- `mfa_required` - (Optional) Require MFA. By default is `false`.

- `mfa_prompt` - (Optional) Prompt for MFA based on the device used, a factor session lifetime, or every sign on attempt: `"DEVICE"`, `"SESSION"` or `"ALWAYS"`.

- `mfa_remember_device` - (Optional) Remember MFA device. The default `false`.

- `mfa_lifetime` - (Optional) Elapsed time before the next MFA challenge.

- `session_idle` - (Optional) Max minutes a session can be idle.",

- `session_lifetime` - (Optional) Max minutes a session is active: Disable = 0.

- `session_persistent` - (Optional) Whether session cookies will last across browser sessions. Okta Administrators can never have persistent session cookies.

- `network_connection` - (Optional) Network selection mode: `"ANYWHERE"`, `"ZONE"`, `"ON_NETWORK"`, or `"OFF_NETWORK"`.

- `network_includes` - (Optional) The network zones to include. Conflicts with `network_excludes`.

- `network_excludes` - (Optional) The network zones to exclude. Conflicts with `network_includes`.

## Attributes Reference

- `id` - ID of the Rule.

- `policyid` - Policy ID.

## Import

A Policy Rule can be imported via the Policy and Rule ID.

```
$ terraform import okta_policy_rule_signon.example <policy id>/<rule id>
```
