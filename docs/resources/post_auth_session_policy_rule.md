---
page_title: "Resource: okta_post_auth_session_policy_rule"
description: |-
  Manages the Post Auth Session Policy Rule.
---

# Resource: okta_post_auth_session_policy_rule

Manages the Post Auth Session Policy Rule. The Post Auth Session Policy has exactly one modifiable rule (non-default). This resource allows you to configure that rule.

~> **IMPORTANT:** This resource cannot be created or deleted, only **imported** and updated. The Post Auth Session Policy rule is pre-provisioned by Okta. You must import the existing rule before managing it with Terraform.

## Import

Before using this resource, you must import the existing rule:

```terraform
terraform import okta_post_auth_session_policy_rule.example <policy_id>/<rule_id>
```

When you run `terraform apply` without importing first, the error message will include the exact import command with the correct policy and rule IDs.

## Example Usage

```terraform
data "okta_post_auth_session_policy" "example" {
}

resource "okta_post_auth_session_policy_rule" "example" {
  policy_id         = data.okta_post_auth_session_policy.example.id
  name              = "Session Protection Rule"
  terminate_session = true
}

```

## Argument Reference

- `policy_id` - (Required) ID of the Post Auth Session Policy. Use the `okta_post_auth_session_policy` data source to get this ID.
- `name` - (Optional) Name of the policy rule.
- `status` - (Optional) Status of the rule: `ACTIVE` or `INACTIVE`. Default is `ACTIVE`.
- `users_excluded` - (Optional) List of user IDs to exclude from this rule.
- `groups_included` - (Optional) List of group IDs to include in this rule.
- `groups_excluded` - (Optional) List of group IDs to exclude from this rule.
- `terminate_session` - (Optional) When true, terminates the user's session when a policy failure is detected. Default is `false`.
- `workflow_id` - (Optional) ID of the Okta Workflow to run when a policy failure is detected.

## Attributes Reference

- `id` - ID of the policy rule.

## Lifecycle

- **Create**: Returns an error with the import command to use
- **Update**: Updates the rule configuration in Okta
- **Delete**: Removes the rule from Terraform state only (the rule remains in Okta)
