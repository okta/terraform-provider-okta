---
page_title: "Data Source: okta_identity_source_groups"
description: |-
  Retrieves an Okta Identity Source Group.
---

# Data Source: okta_identity_source_groups

Retrieves a group record from an Okta Identity Source by external ID.

## Example Usage

```terraform
data "okta_identity_source_groups" "example" {
  identity_source_id = "<identity-source-id>"
  external_id        = "GROUPEXT123456EXAMPLE"
}

output "group_okta_id" {
  value = data.okta_identity_source_groups.example.id
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source.
- `external_id` (String) The external ID of the group in the identity source.

### Optional

- `id` (String) Okta-assigned ID of the group. May be provided to pre-populate or verify the Okta ID.

### Read-Only

- `id` (String) Okta-assigned ID of the group (populated after read).
