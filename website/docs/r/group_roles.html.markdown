---
layout: 'okta'
page_title: 'Okta: okta_group_roles'
sidebar_current: 'docs-okta-resource-group-roles'
description: |-
  Creates Group level Admin Role Assignments.
---

# okta_group_roles

Creates Group level Admin Role Assignments.

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

- `admin_roles` - (Required) Admin roles associated with the group. It can be any of the following values `"SUPER_ADMIN"`, `"ORG_ADMIN"`, `"APP_ADMIN"`, `"USER_ADMIN"`, `"HELP_DESK_ADMIN"`, `"READ_ONLY_ADMIN"`, `"MOBILE_ADMIN"`, `"API_ACCESS_MANAGEMENT_ADMIN"`, `"REPORT_ADMIN"`, `"GROUP_MEMBERSHIP_ADMIN"`.

## Attributes Reference

- `id` - The ID of the Group Role Assignment.

## Import

Group Role Assignment can be imported via the Okta Group ID.

```
$ terraform import okta_group_roles.example <group id>
```
