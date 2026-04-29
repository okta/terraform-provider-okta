---
page_title: "Resource: okta_identity_source_bulk_groups_delete"
description: |-
  Uploads a list of group external IDs to an Okta Identity Source session for deletion.
---

# Resource: okta_identity_source_bulk_groups_delete

Uploads a list of group external IDs into an Okta Identity Source session, staging them for deletion from Okta. Use this resource alongside `okta_identity_source_session`. The operation is write-only: Okta does not expose a GET endpoint for the uploaded deletion list, so state is preserved from the create call.

~> **Note:** This resource does not support deletion. Removing it from configuration will emit a warning but take no action against Okta.

## Example Usage

```terraform
resource "okta_identity_source_session" "example" {
  identity_source_id = "<identity-source-id>"
}

resource "okta_identity_source_bulk_groups_delete" "example" {
  identity_source_id = okta_identity_source_session.example.identity_source_id
  session_id         = okta_identity_source_session.example.id

  external_ids = ["GROUPEXT123456EXAMPLE"]
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source. Forces replacement when changed.
- `session_id` (String) ID of the identity source session to upload deletion data into. Forces replacement when changed.

### Optional

- `external_ids` (List of String) Array of external IDs of groups to delete in Okta.

### Read-Only

- `id` (String) The unique identifier for the resource (set to the session ID).

## Import

Import using `{identity_source_id}/{session_id}/{id}`:

```shell
terraform import okta_identity_source_bulk_groups_delete.example <identity-source-id>/<session-id>/<id>
```
