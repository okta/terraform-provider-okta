---
page_title: "okta_identity_source_bulk_groups_delete Resource - terraform-provider-okta"
description: |-
  Manages an Okta Identity Source Bulk Groups Delete resource.
---

# okta_identity_source_bulk_groups_delete (Resource)

Manages an Okta Identity Source Bulk Groups Delete.

## Example Usage

```hcl
resource "okta_identity_source_bulk_groups_delete" "example" {
  identity_source_id = "<identity-source-id>"
  session_id = "<session-id>"

  # Optional fields
  # external_ids = "<external_ids>"
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source
- `session_id` (String) ID of the identity source session

### Optional

- `external_ids` (List) Array of external IDs of groups that need to be deleted in Okta

### Read-Only

- `id` (String) The unique identifier for the resource.

## Import

Import using `{identity_source_id}/{session_id}/{id}`:

```shell
terraform import okta_identity_source_bulk_groups_delete.example <identity_source_id> <session_id> <id>
```
