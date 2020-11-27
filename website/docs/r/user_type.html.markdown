---
layout: 'okta'
page_title: 'Okta: okta_user_type'
sidebar_current: 'docs-okta-resource-user-type'
description: |-
  Creates a User Type.
---

# okta_user_type

Creates a User type.

This resource allows you to create and configure a User Type.

## Example Usage

```hcl
resource "okta_user_type" "example" {
  name   = "example"
  display_name = "example"
  description = "example"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the User Type.

- `display_name` - (Required) Display Name of the User Type.

- `description` - (Optional) Description of the User Type.

## Attributes Reference

- `id` - The ID of the User Type.

## Import

A User Type can be imported via the Okta ID.

```
$ terraform import okta_user_type.example <user type id>
```
