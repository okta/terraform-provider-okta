---
layout: 'okta' page_title: 'Okta: okta_group_role_group_targets' sidebar_current: 'docs-okta-resource-group-role'
description: |- Manages target groups for admin role assignments.
---

# okta_group_role_group_targets

Manages target groups for admin role assignments.

This resource allows you to assign a list of target groups to group admin role assignments, when the role type supports
it. This will generally mean the user with this administrative ability can only perform the act on the groups or members
of those groups.

This resource only supports targeting admin role assignments of the following types:

- `GROUP_MEMBERSHIP_ADMIN`
- `HELP_DESK_ADMIN`
- `USER_ADMIN`

## Example Usage

```hcl
resource "okta_group_role_group_targets" "example" {
  group_id = "<group id>"
  role_id  = "<role id>"
  target_group_list = [
    "target_group_id",
    "target_group_id"
  ]
}
```

## Argument Reference

The following arguments are supported:

- `group_id` - (Required) The ID of group with the admin role being modified.

- `role_id` - (Required) Admin role ID you wish to target.

- `target_group_list` - (Required) A list of group IDs you would like as the targets of the admin role.

## Attributes Reference

N/A

## Import

N/A
