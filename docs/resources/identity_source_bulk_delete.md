---
page_title: "Resource: okta_identity_source_bulk_delete"
description: |-
  Uploads a batch of user external IDs to an Okta Identity Source session for deletion.
---

# Resource: okta_identity_source_bulk_delete

Uploads a list of user external IDs into an Okta Identity Source session, staging them for deletion from Okta. Use this resource alongside `okta_identity_source_session`. The operation is write-only: Okta does not expose a GET endpoint for the uploaded deletion list, so state is preserved from the create call.

~> **Note:** This resource does not support deletion. Removing it from configuration will emit a warning but take no action against Okta.

## Example Usage

```terraform
resource "okta_identity_source_session" "example" {
  identity_source_id = "<identity-source-id>"
}

resource "okta_identity_source_bulk_delete" "example" {
  identity_source_id = okta_identity_source_session.example.identity_source_id
  session_id         = okta_identity_source_session.example.id
  entity_type        = "USERS"

  profiles {
    external_id = "USEREXT123456EXAMPLE"
  }
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source. Forces replacement when changed.
- `session_id` (String) ID of the identity source session to upload deletion data into. Forces replacement when changed.

### Optional

- `entity_type` (String) The type of data to bulk delete in a session. Currently, only `USERS` is supported.
- `profiles` (Block List, Optional) Array of profiles identifying entities to delete. (see [below for nested schema](#nested-schema-for-profiles))

### Read-Only

- `id` (String) The unique identifier for the resource (set to the session ID).

### Nested Schema for `profiles`

Optional:

- `external_id` (String) The external ID of the entity to delete in Okta.

## Import

Import using `{identity_source_id}/{session_id}/{id}`:

```shell
terraform import okta_identity_source_bulk_delete.example <identity-source-id>/<session-id>/<id>
```
