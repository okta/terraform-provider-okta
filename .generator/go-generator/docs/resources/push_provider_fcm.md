---
page_title: "okta_push_provider_fcm Resource - terraform-provider-okta"
description: |-
  Manages an Okta Push Provider Fcm resource.
---

# okta_push_provider_fcm (Resource)

Manages an Okta Push Provider Fcm.

## Example Usage

```hcl
resource "okta_push_provider_fcm" "example" {
  provider_type = "<provider_type>"

  # Optional fields
  # configuration = "<configuration>"
  # name = "Example Name"
}
```

## Schema

### Required

- `provider_type` (String) Discriminator field identifying the variant type. Must be set to \

### Optional

- `configuration` (String) Configuration
- `name` (String) Display name of the push provider

### Read-Only

- `id` (String) The unique identifier for the resource.
- `last_updated_date` (String) Timestamp when the Push Provider was last modified
