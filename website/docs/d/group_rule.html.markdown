---
layout: "okta"
page_title: "Okta: okta_group_rule"
sidebar_current: "docs-okta-datasource-group-rule"
description: |- Get a group rule from Okta.
---

# okta_group_rule

Use this data source to retrieve a group rule from Okta.

## Example Usage

```hcl
data "okta_group_rule" "test" {
  id = okta_group_rule.example.id
}
```

## Arguments Reference

- `id` - (Required) ID of the group rule to retrieve.


## Attributes Reference

- `id` - The ID of the Group Rule.

- `name` - The name of the Group Rule.

- `group_assignments` - The list of group ids to assign the users to.

- `expression_type` - The expression type to use to invoke the rule.

- `expression_value` - The expression value.

- `status` - The status of the group rule.

- `users_excluded` - The list of user IDs that would be excluded when rules are processed.
