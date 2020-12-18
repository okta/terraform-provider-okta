---
layout: 'okta'
page_title: 'Okta: okta_policy_rule_mfa'
sidebar_current: 'docs-okta-resource-policy-rule-mfa'
description: |-
  Creates an MFA Policy Rule.
---

# okta_policy_rule_mfa

Creates an MFA Policy Rule.

This resource allows you to create and configure an MFA Policy Rule.

## Argument Reference

The following arguments are supported:

- `policyid` - (Required) Policy ID.

- `name` - (Required) Policy Rule Name.

- `priority` - (Optional) Policy Rule Priority, this attribute can be set to a valid priority. To avoid endless diff situation we error if an invalid priority is provided. API defaults it to the last/lowest if not there.

- `status` - (Optional) Policy Rule Status: `"ACTIVE"` or `"INACTIVE"`.

- `enroll` - (Optional) When a user should be prompted for MFA. It can be `"CHALLENGE"`, `"LOGIN"`, or `"NEVER"`.

- `network_connection` - (Optional) Network selection mode: `"ANYWHERE"`, `"ZONE"`, `"ON_NETWORK"`, or `"OFF_NETWORK"`.

- `network_includes` - (Optional) The network zones to include. Conflicts with `network_excludes`.

- `network_excludes` - (Optional) The network zones to exclude. Conflicts with `network_includes`.

## Attributes Reference

- `id` - ID of the Rule.

- `policyid` - Policy ID.

## Import

A Policy Rule can be imported via the Policy and Rule ID.

```
$ terraform import okta_policy_rule_mfa.example <policy id>/<rule id>
```
