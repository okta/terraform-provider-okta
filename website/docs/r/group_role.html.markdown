---
layout: 'okta' page_title: 'Okta: okta_group_role' sidebar_current: 'docs-okta-resource-group-role' description: |-
Assigns Admin roles to Okta Groups.
---

# okta_group_role

Assigns Admin roles to Okta Groups.

This resource allows you to assign Okta administrator roles to Okta Groups. This resource provides a one-to-one
interface between the Okta group and the admin role.

## Example Usage

```hcl
resource "okta_group_role" "example" {
  group_id  = "<group id>"
  role_type = "READ_ONLY_ADMIN"
}
```

## Argument Reference

The following arguments are supported:

- `group_id` - (Required) The ID of group to attach admin roles to.

- `role_type` - (Required) Admin role assigned to the group. It can be any one of the following values `"SUPER_ADMIN"`
  , `"ORG_ADMIN"`, `"APP_ADMIN"`, `"USER_ADMIN"`, `"HELP_DESK_ADMIN"`, `"READ_ONLY_ADMIN"`
  , `"MOBILE_ADMIN"`, `"API_ACCESS_MANAGEMENT_ADMIN"`, `"REPORT_ADMIN"`, `"GROUP_MEMBERSHIP_ADMIN"`.

- `target_group_list` - (Optional) A list of group IDs you would like as the targets of the admin role.
    - Only supported when used with the role types: `GROUP_MEMBERSHIP_ADMIN`, `HELP_DESK_ADMIN`, or `USER_ADMIN`.

- `target_app_list` - (Optional) A list of app names (name represents set of app instances, like 'salesforce' or '
  facebook'), or a combination of app name and app instance ID (like 'facebook.0oapsqQ6dv19pqyEo0g3') you would like as
  the targets of the admin role.
    - Only supported when used with the role type `"APP_ADMIN"`.

## Attributes Reference

- `id` - The ID of the Group Role Assignment.

## Import

Individual admin role assignment can be imported by passing the group and role assignment IDs as follows:

```
$ terraform import okta_group_role.example <group id>/<role id>
```
