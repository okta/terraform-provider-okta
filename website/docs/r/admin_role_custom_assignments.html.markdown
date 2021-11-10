---
layout: 'okta'
page_title: 'Okta: okta_admin_role_custom_assignments'
sidebar_current: 'docs-okta-resource-okta-admin-role-custom-assignments'
description: |-
    Manages custom roles assignments
---

# okta_admin_role_custom_assignments

This resource allows the assignment and unassignment of Custom Roles. The `members` field supports these type of resources:
 - Groups
 - Users

~> **NOTE:** This an `Early Access` feature.

## Example Usage

```hcl
locals {
  org_url = "https://mycompany.okta.com"
}

resource "okta_admin_role_custom" "test" {
  label       = "SomeUsersAndApps"
  description = "Manage apps assignments and users"
  permissions = ["okta.apps.assignment.manage", "okta.users.manage", "okta.apps.manage"]
}

resource "okta_resource_set" "test" {
  label       = "UsersWithApp"
  description = "All the users and SWA app"
  resources   = [
    format("%s/api/v1/users", local.org_url),
    format("%s/api/v1/apps/%s", local.org_url, okta_app_swa.test.id)
  ]
}

// this user and group will manage the set of resources based on the permissions specified in the custom role
resource "okta_admin_role_custom_assignments" "test" {
  resource_set_id = okta_resource_set.test.id
  custom_role_id  = okta_admin_role_custom.test.id
  members         = [
    format("%s/api/v1/users/%s", local.org_url, okta_user.test.id),
    format("%s/api/v1/groups/%s", local.org_url, okta_group.test.id)
  ]
}

// this user will have `CUSTOM` role assigned, but it won't appear in the `admin_roles` for that user,
// since direct assignment of custom roles is not allowed
resource "okta_user" "test" {
  first_name = "Paul"
  last_name  = "Atreides"
  login      = "no-reply@caladan.planet"
  email      = "no-reply@caladan.planet"
}

resource "okta_app_swa" "test" {
  label          = "My SWA App"
  button_field   = "btn-login"
  password_field = "txtbox-password"
  username_field = "txtbox-username"
  url            = "https://example.com/login.html"
}

resource "okta_group" "test" {
  name        = "General"
  description = "General Group"
}
```

## Argument Reference

The following arguments are supported:

- `resource_set_id` - (Required) ID of the target Resource Set.

- `custom_role_id` - (Required) ID of the Custom Role.

- `members` - (Optional) The hrefs that point to User(s) and/or Group(s) that receive the Role. At least one
  permission must be specified when creating custom role.

## Attributes Reference

- `id` - ID of this resource in `resource_set_id/custom_role_id` format.

## Import

Okta Custom Admin Role Assignments can be imported via the Okta ID.

```
$ terraform import okta_admin_role_custom_assignments.example <resource_set_id>/<custom_role_id>
```
