---
page_title: "Resource: okta_session_violation_policy_rule"
description: |-
  Manages the Session Violation Detection Policy Rule.
---

# Resource: okta_session_violation_policy_rule

Manages the Session Violation Detection Policy Rule. The Session Violation Detection Policy has exactly one modifiable rule (non-default). This resource allows you to configure that rule.

~> **IMPORTANT:** This resource cannot be created or deleted, only **imported** and updated. The Session Violation Detection Policy rule is pre-provisioned by Okta. You must import the existing rule before managing it with Terraform.

## Import

Before using this resource, you must import the existing rule:

{{codefile "shell" "examples/resources/okta_session_violation_policy_rule/import.sh"}}

Use the `okta_session_violation_policy` data source to retrieve the `policy_id` and `rule_id` needed for the import command.

## Example Usage

```terraform
data "okta_session_violation_policy" "example" {
}

resource "okta_session_violation_policy_rule" "example" {
  policy_id                 = data.okta_session_violation_policy.example.id
  name                      = "Session Violation Rule"
  min_risk_level            = "HIGH"
  policy_evaluation_enabled = true
}
```

## Argument Reference

- `policy_id` - (Required) ID of the Session Violation Detection Policy. Use the `okta_session_violation_policy` data source to get this ID.
- `min_risk_level` - (Required) The minimum risk level that triggers the rule. Valid values: `LOW`, `MEDIUM`, `HIGH`.
- `name` - (Optional) Name of the policy rule.
- `status` - (Optional) Status of the rule: `ACTIVE` or `INACTIVE`. Default is `ACTIVE`.
- `network_connection` - (Optional) Network selection mode. Valid values: `ANYWHERE`, `ZONE`, `ON_NETWORK`, `OFF_NETWORK`.
- `network_includes` - (Optional) List of network zone IDs to include. Required when `network_connection` is set to `ZONE`.
- `network_excludes` - (Optional) List of network zone IDs to exclude. Required when `network_connection` is set to `ZONE`.
- `policy_evaluation_enabled` - (Optional) When `true`, the sign-on policies of the session are evaluated when a session violation is detected. Default is `true`.
- `priority` - (Optional) Priority of the rule. Rules are evaluated in priority order.

## Attributes Reference

- `id` - ID of the policy rule.

## Lifecycle

- **Create**: Returns an error with the import command to use
- **Update**: Updates the rule configuration in Okta
- **Delete**: Removes the rule from Terraform state only (the rule remains in Okta)
