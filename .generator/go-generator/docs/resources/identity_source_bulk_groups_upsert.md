---
page_title: "okta_identity_source_bulk_groups_upsert Resource - terraform-provider-okta"
description: |-
  Manages an Okta Identity Source Bulk Groups Upsert resource.
---

# okta_identity_source_bulk_groups_upsert (Resource)

Manages an Okta Identity Source Bulk Groups Upsert.

## Example Usage

```hcl
resource "okta_identity_source_bulk_groups_upsert" "example" {
  identity_source_id = "<identity-source-id>"
  session_id = "<session-id>"

  # Optional fields
  # external_id = "<external-id>"
  # description = "Example description"
  # display_name = "Example Display Name"
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source
- `session_id` (String) ID of the identity source session

### Optional

- `external_id` (String) The external ID of the group that needs to be created or updated in Okta
- `description` (String) Description of the group
- `display_name` (String) Name of the group

### Read-Only

- `id` (String) The unique identifier for the resource.

## Import

Import using `{identity_source_id}/{session_id}/{id}`:

```shell
terraform import okta_identity_source_bulk_groups_upsert.example <identity_source_id> <session_id> <id>
```
