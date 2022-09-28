---
layout: 'okta'
page_title: 'Okta: okta_group_roles'
sidebar_current: 'docs-okta-resource-group-roles'
description: |-
  Creates Group level Admin Role Assignments.
---

# okta_group_roles

~> **DEPRECATED:** This resource is deprecated and will be removed in favor of using `okta_group_role`, please migrate as soon as possible

This resource allows you to create and configure Group level Admin Role Assignments.

## Example Usage

```hcl
resource "okta_group_roles" "example" {
  group_id    = "<group id>"
  admin_roles = ["SUPER_ADMIN"]
}
```

## Argument Reference

The following arguments are supported:

- `group_id` - (Required) The ID of group to attach admin roles to.

- `admin_roles` - (Required) Admin roles associated with the group. It can be any of the following values:
  `"API_ADMIN"`,
  `"APP_ADMIN"`,
  `"CUSTOM"`,
  `"GROUP_MEMBERSHIP_ADMIN"`,
  `"HELP_DESK_ADMIN"`,
  `"MOBILE_ADMIN"`,
  `"ORG_ADMIN"`,
  `"READ_ONLY_ADMIN"`,
  `"REPORT_ADMIN"`,
  `"SUPER_ADMIN"`,
  `"USER_ADMIN"`
  .

## Attributes Reference

- `id` - The ID of the Group Role Assignment.

## Import

Group Role Assignment can be imported via the Okta Group ID.

```
$ terraform import okta_group_roles.example &#60;group id&#62;
```
