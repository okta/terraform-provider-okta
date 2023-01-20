---
layout: 'okta'
page_title: 'Okta: okta_group_role'
sidebar_current: 'docs-okta-resource-group-role'
description: |-
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

- `role_type` - (Required) Admin role assigned to the group. It can be any one of the following values:
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
  . See [API Docs](https://developer.okta.com/docs/reference/api/roles/#role-types).


  - `"USER_ADMIN"` is the Group Administrator.


- `target_group_list` - (Optional) A list of group IDs you would like as the targets of the admin role.
    - Only supported when used with the role types: `GROUP_MEMBERSHIP_ADMIN`, `HELP_DESK_ADMIN`, or `USER_ADMIN`.

- `target_app_list` - (Optional) A list of app names (name represents set of app instances, like 'salesforce' or '
  facebook'), or a combination of app name and app instance ID (like 'facebook.0oapsqQ6dv19pqyEo0g3') you would like as
  the targets of the admin role.
    - Only supported when used with the role type `"APP_ADMIN"`.

- `disable_notifications` - (Optional) When this setting is enabled, the admins won't receive any of the default Okta
  administrator emails. These admins also won't have access to contact Okta Support and open support cases on behalf of your org.

## Attributes Reference

- `id` - The ID of the Group Role Assignment.

## Import

Individual admin role assignment can be imported by passing the group and role assignment IDs as follows:

```
$ terraform import okta_group_role.example &#60;group id&#62;/&#60;role id&#62;
```
