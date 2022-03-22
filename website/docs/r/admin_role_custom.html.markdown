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
  permission must be specified when creating custom role. Valid values: `"okta.apps.assignment.manage"`,
`"okta.apps.manage"`,
`"okta.apps.read"`,
`"okta.groups.appAssignment.manage"`,
`"okta.groups.create"`,
`"okta.groups.manage"`,
`"okta.groups.members.manage"`,
`"okta.groups.read"`,
`"okta.profilesource.import.run"`,
`"okta.users.appAssignment.manage"`,
`"okta.users.create"`,
`"okta.users.credentials.expirePassword"`,
`"okta.users.credentials.manage"`,
`"okta.users.credentials.resetFactors"`,
`"okta.users.credentials.resetPassword"`,
`"okta.users.groupMembership.manage"`,
`"okta.users.lifecycle.activate"`,
`"okta.users.lifecycle.clearSessions"`,
`"okta.users.lifecycle.deactivate"`,
`"okta.users.lifecycle.delete"`,
`"okta.users.lifecycle.manage"`,
`"okta.users.lifecycle.suspend"`,
`"okta.users.lifecycle.unlock"`,
`"okta.users.lifecycle.unsuspend"`,
`"okta.users.manage"`,
`"okta.users.read"`,
`"okta.users.userprofile.manage"`.

## Attributes Reference

- `id` - Custom Role ID

## Import

Okta Custom Admin Role can be imported via the Okta ID.

```
$ terraform import okta_admin_role_custom.example <custom role id>
```
