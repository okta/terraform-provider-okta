---
page_title: "okta_user_type Resource - terraform-provider-okta"
description: |-
  Manages an Okta User Type resource.
---

# okta_user_type (Resource)

Manages an Okta User Type.

## Example Usage

```hcl
resource "okta_user_type" "example" {
  display_name = "Example Display Name"
  name = "Example Name"

  # Optional fields
  # description = "Example description"
}
```

## Schema

### Required

- `display_name` (String) The human-readable name of the user type
- `name` (String) The name of the user type. The name must start with A-Z or a-z and contain only A-Z, a-z, 0-9, or underscore (_) characters. This value becomes read-only after creation and can't be updated.

### Optional

- `description` (String) The human-readable description of the user type

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) A timestamp from when the user type was created
- `created_by` (String) The user ID of the account that created the user type
- `default` (Boolean) A boolean value to indicate if this is the default user type
- `last_updated` (String) A timestamp from when the user type was most recently updated
- `last_updated_by` (String) The user ID of the most recent account to edit the user type
