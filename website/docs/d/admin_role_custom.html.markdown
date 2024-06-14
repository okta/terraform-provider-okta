---
layout: 'okta'
page_title: 'Okta: okta_admin_role_custom'
sidebar_current: 'docs-okta-datasource-okta-admin-role-custom'
description: |-
    Get a custom admin role from Okta.
---

# okta_admin_role_custom

Use this data source to retrieve a custom admin role from Okta.

## Example Usage

```hcl
data "okta_admin_role_custom" "example" {
  label       = "ExampleRole"
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Optional) The ID of the custom role to retrieve, conflicts with `label`.

- `label` - (Optional) The name of the custom role to retrieve, conflicts with `id`.

## Attributes Reference

- `id` - The ID of the custom role.

- `label` - The name of the custom role.

- `description` - The human-readable description of the custom role.

- `permissions` - The list of permissions that the custom role grants.
