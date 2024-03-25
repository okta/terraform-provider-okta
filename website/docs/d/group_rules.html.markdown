---
layout: "okta"
page_title: "Okta: okta_group_rules"
sidebar_current: "docs-okta-datasource-group-rules"
description: |- Get a list of group rules from Okta.
---

# okta_group_rules

Use this data source to retrieve a list of group rules from Okta.

## Example Usage

```hcl
data "okta_group_rules" "department_rules" {
  name_prefix = "Rule-Department-"
}
```

## Arguments Reference

- `name_prefix` - (Optional) The name prefix of the target Group Rule(s).

## Attributes Reference

- `rules` - collection of group rules retrieved from Okta with the following properties.

  - `id` - The ID of the Group Rule.

  - `name` - The name of the Group Rule.

  - `group_assignments` - The list of group ids to assign the users to.

  - `expression_type` - The expression type to use to invoke the rule.

  - `expression_value` - The expression value.

  - `status` - The status of the group rule.

  - `users_excluded` - The list of user IDs that would be excluded when rules are processed.
