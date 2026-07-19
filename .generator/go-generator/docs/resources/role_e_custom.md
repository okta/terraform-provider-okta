---
page_title: "okta_role_e_custom Resource - terraform-provider-okta"
description: |-
  Manages an Okta Role E Custom resource.
---

# okta_role_e_custom (Resource)

Manages an Okta Role E Custom.

## Example Usage

```hcl
resource "okta_role_e_custom" "example" {
  description = "Example description"
  label = "Example Label"
  permissions = "<permissions>"
}
```

## Schema

### Required

- `description` (String) Description of the role
- `label` (String) Unique label for the role
- `permissions` (List) Array of permissions that the role grants. See [Permissions](/openapi/okta-management/guides/permissions).

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the role was created
- `last_updated` (String) Timestamp when the role was last updated
