---
page_title: "okta_identity_source_bulk_delete Resource - terraform-provider-okta"
description: |-
  Manages an Okta Identity Source Bulk Delete resource.
---

# okta_identity_source_bulk_delete (Resource)

Manages an Okta Identity Source Bulk Delete.

## Example Usage

```hcl
resource "okta_identity_source_bulk_delete" "example" {
  identity_source_id = "<identity-source-id>"
  session_id = "<session-id>"

  # Optional fields
  # entity_type = "<entity_type>"
  # external_id = "<external-id>"
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source
- `session_id` (String) ID of the identity source session

### Optional

- `entity_type` (String) The type of data to bulk delete in a session. Currently, only `USERS` is supported.
- `external_id` (String) The external ID of the entity that needs to be deleted in Okta

### Read-Only

- `id` (String) The unique identifier for the resource.

## Import

Import using `{identity_source_id}/{session_id}/{id}`:

```shell
terraform import okta_identity_source_bulk_delete.example <identity_source_id> <session_id> <id>
```
