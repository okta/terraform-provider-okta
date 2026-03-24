---
page_title: "Resource: okta_entity_risk_policy_rule"
description: |-
  Manages an Entity Risk Policy Rule in Okta.
---

# Resource: okta_entity_risk_policy_rule

Manages an Entity Risk Policy Rule. Entity Risk Policy rules define automated responses to identity threats detected by Okta's Identity Threat Protection (ITP).

~> **NOTE:** Entity Risk Policy is automatically created when Identity Threat Protection (ITP) is enabled. Use the `okta_entity_risk_policy` data source to get the policy ID. The default policy rule (priority 99) cannot be imported or modified.

## Example Usage

### Basic Rule - Terminate Sessions on High Risk

```terraform
data "okta_entity_risk_policy" "example" {
}

resource "okta_entity_risk_policy_rule" "high_risk" {
  policy_id              = data.okta_entity_risk_policy.example.id
  name                   = "High Risk - Terminate Sessions"
  risk_level             = "HIGH"
  terminate_all_sessions = true
}
```

### Rule with Group Targeting

```terraform
data "okta_entity_risk_policy" "example" {
}

data "okta_group" "privileged_users" {
  name = "Privileged Users"
}

resource "okta_entity_risk_policy_rule" "privileged_high_risk" {
  policy_id              = data.okta_entity_risk_policy.example.id
  name                   = "Privileged Users - High Risk"
  risk_level             = "HIGH"
  terminate_all_sessions = true
  groups_included        = [data.okta_group.privileged_users.id]
}
```

### Rule with Workflow Integration

```terraform
data "okta_entity_risk_policy" "example" {
}

resource "okta_entity_risk_policy_rule" "workflow_rule" {
  policy_id   = data.okta_entity_risk_policy.example.id
  name        = "Low Risk - Run Workflow"
  risk_level  = "LOW"
  workflow_id = "your-workflow-id"
}
```

## Argument Reference

The following arguments are supported:

- `policy_id` - (Required) ID of the Entity Risk Policy. Use the `okta_entity_risk_policy` data source to get this ID.
- `name` - (Required) Name of the policy rule.
- `risk_level` - (Required) Risk level to match. Valid values: `HIGH`, `MEDIUM`, `LOW`, `ANY`.
- `status` - (Optional) Status of the rule. Valid values: `ACTIVE`, `INACTIVE`. Default: `ACTIVE`.
- `priority` - (Optional) Priority of the rule. Rules are evaluated in priority order.
- `users_excluded` - (Optional) List of user IDs to exclude from this rule.
- `groups_included` - (Optional) List of group IDs to include in this rule.
- `groups_excluded` - (Optional) List of group IDs to exclude from this rule.
- `terminate_all_sessions` - (Optional) When true, terminates all active sessions for the user when a risk event is detected. Default: `false`.
- `workflow_id` - (Optional) ID of the Okta Workflow to run when a risk event is detected.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the policy rule.

## Import

Entity Risk Policy Rules can be imported using the format `<policy_id>/<rule_id>`:

```shell
terraform import okta_entity_risk_policy_rule.example 00p1234567890/0pr1234567890
```
