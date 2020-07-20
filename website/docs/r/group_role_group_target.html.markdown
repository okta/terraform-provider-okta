---
layout: "okta"
page_title: "Okta: okta_group_role_group_target"
sidebar_current: "docs-okta-resource-group-role-group-target"
description: |-
  Creates Group targets for Group level Admin Role Assignments.
---

# okta_group_role

Creates Group targets for Group level Admin Role Assignments.

This resource allows you to create and configure Group targets for Group level Admin Role Assignments (i.e., Group Admin assignments).

## Example Usage

```hcl
resource "okta_group_role" "example" {
  group_id    = "<group id>"
  role_type   = "USER_ADMIN"
}

resource "okta_group_role_group_target" "example_target" {
  group_id          = "<group id>"
  role_id           = okta_group_role.example.id
  target_group_id   = "<target group id>"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required) The ID of the group to attach the target to.

* `role_id` - (Required) The ID of the role to attach the target to.

* `target_group_id` - (Required) The ID of the target group.

## Attributes Reference

* `id` - The ID of the Group Role Assignment.


