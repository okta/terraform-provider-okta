---
layout: 'okta'
page_title: 'Okta: okta_group_rule'
sidebar_current: 'docs-okta-resource-group-rule'
description: |-
  Creates an Okta Group Rule.
---

# okta_group_rule

Creates an Okta Group Rule.

This resource allows you to create and configure an Okta Group Rule.

## Example Usage

```hcl
resource "okta_group_rule" "example" {
  name              = "example"
  status            = "ACTIVE"
  group_assignments = ["<group id>"]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"andy\")"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the Group Rule.

- `group_assignments` - (Required) The list of group ids to assign the users to.

- `expression_type` - (Optional) The expression type to use to invoke the rule. The default is `"urn:okta:expression:1.0"`.

- `expression_value` - (Required) The expression value.

- `status` - (Optional) The status of the group rule.

## Attributes Reference

- `id` - The ID of the Group Rule.

## Import

An Okta Group Rule can be imported via the Okta ID.

```
$ terraform import okta_group_rule.example <group rule id>
```
