---
page_title: "okta_hook_key Resource - terraform-provider-okta"
description: |-
  Manages an Okta Hook Key resource.
---

# okta_hook_key (Resource)

Manages an Okta Hook Key.

## Example Usage

```hcl
resource "okta_hook_key" "example" {

  # Optional fields
  # name = "Example Name"
}
```

## Schema

### Optional

- `name` (String) Display name of the key

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the key was created
- `is_used` (String) Whether this key is currently in use by other applications
- `key_id` (String) The alias of the public key
- `last_updated` (String) Timestamp when the key was updated
