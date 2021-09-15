---
layout: 'okta'
page_title: 'Okta: okta_group'
sidebar_current: 'docs-okta-resource-group'
description: |-
  Creates an Okta Group.
---

# okta_group

Creates an Okta Group.

This resource allows you to create and configure an Okta Group.

## Example Usage

```hcl
resource "okta_group" "example" {
  name        = "Example"
  description = "My Example Group"
}
```

Ignore users sync
```hcl
resource "okta_group" "example_skip" {
  name        = "Example"
  description = "My Example Group"
  skip_users  = true
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the Okta Group.

- `description` - (Optional) The description of the Okta Group.

- `users` - (Optional) The users associated with the group. This can also be done per user.
`DEPRECATED`: Please replace usage with the `okta_group_memberships` resource.

- `skip_users` - (Optional) Indicator that allows the app to skip `users` sync (it's also can be provided during import). Default is `false`.

## Attributes Reference

- `id` - The ID of the Okta Group.

## Import

An Okta Group can be imported via the Okta ID.

```
$ terraform import okta_group.example <group id>
```

It's also possible to import group without users. In this case ID will look like this:

```
$ terraform import okta_app_basic_auth.example <group id>/skip_users
```
