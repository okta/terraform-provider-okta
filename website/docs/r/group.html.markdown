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

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the Okta Group.

- `description` - (Optional) The description of the Okta Group.

- `users` - (Optional) The users associated with the group. This can also be done per user.

## Attributes Reference

- `id` - The ID of the Okta Group.

## Import

An Okta Group can be imported via the Okta ID.

```
$ terraform import okta_group.example <group id>
```
