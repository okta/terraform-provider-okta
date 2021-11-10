---
layout: 'okta'
page_title: 'Okta: okta_admin_role_custom'
sidebar_current: 'docs-okta-resource-okta-admin-role-custom'
description: |-
    Manages custom roles.
---

# okta_admin_role_custom

These operations allow the creation and manipulation of custom roles as custom collections of permissions.

~> **NOTE:** This an `Early Access` feature.

## Example Usage

```hcl
resource "okta_admin_role_custom" "example" {
  label       = "AppAssignmentManager"
  description = "This role allows app assignment management"
  permissions = ["okta.apps.assignment.manage"]
}
```

## Argument Reference

The following arguments are supported:

- `label` - (Required) The name given to the new Role.

- `description` - (Required) A human-readable description of the new Role.

- `permissions` - (Optional) The permissions that the new Role grants. At least one
  permission must be specified when creating custom role. Valid values:`"okta.users.manage"`, 
 `"okta.users.create"`,`"okta.users.read"`,`"okta.users.credentials.manage"`,`"okta.users.userprofile.manage"`, 
 `"okta.users.lifecycle.manage"`,`"okta.users.groupMembership.manage"`,`"okta.users.appAssignment.manage"`,
 `"okta.groups.manage"`,`"okta.groups.create"`,`"okta.groups.members.manage"`,`"okta.groups.read"`,
 `"okta.groups.appAssignment.manage"`,`"okta.apps.read"`,`"okta.apps.manage"`,`"okta.apps.assignment.manage"`. 

## Attributes Reference

- `id` - Custom Role ID

## Import

Okta Custom Admin Role can be imported via the Okta ID.

```
$ terraform import okta_admin_role_custom.example <custom role id>
```
