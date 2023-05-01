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

Custom profile attributes
```hcl
resource "okta_group" "example" {
  name        = "Example"
  description = "My Example Group"
  custom_profile_attributes = jsonencode({
    "example1" = "testing1234",
    "example2" = true,
    "example3" = 54321
  })
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the Okta Group.

- `description` - (Optional) The description of the Okta Group.

- `custom_profile_attributes` - (Optional) raw JSON containing all custom profile attributes.

## Attributes Reference

- `id` - The ID of the Okta Group.

## Import

An Okta Group can be imported via the Okta ID.

```
$ terraform import okta_group.example &#60;group id&#62;
```