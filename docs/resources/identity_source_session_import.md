---
page_title: "Resource: okta_identity_source_session_import"
description: |-
  Triggers the import of staged identity source data into Okta.
---

# Resource: okta_identity_source_session_import

Triggers the import of staged identity source data into Okta by calling the `startImportFromIdentitySource` API. This resource must be created **after** all bulk operation resources (`okta_identity_source_bulk_upsert`, `okta_identity_source_bulk_delete`, `okta_identity_source_bulk_groups_upsert`, `okta_identity_source_bulk_groups_delete`) for the same session have been applied.

~> **Note:** This resource does not support deletion. Removing it from configuration will emit a warning but does not undo the import in Okta.

## Example Usage

```terraform
resource "okta_identity_source_session" "example" {
  identity_source_id = "<identity-source-id>"
}

resource "okta_identity_source_bulk_upsert" "example" {
  identity_source_id = okta_identity_source_session.example.identity_source_id
  session_id         = okta_identity_source_session.example.id

  entity_type = "USERS"

  profiles {
    external_id = "USEREXT123456EXAMPLE"

    profile {
      user_name  = "user@example.com"
      first_name = "Jane"
      last_name  = "Doe"
      email      = "user@example.com"
    }
  }
}

resource "okta_identity_source_session_import" "example" {
  identity_source_id = okta_identity_source_session.example.identity_source_id
  session_id         = okta_identity_source_session.example.id

  depends_on = [okta_identity_source_bulk_upsert.example]
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source. Forces replacement when changed.
- `session_id` (String) ID of the identity source session to trigger import for. Forces replacement when changed.

### Read-Only

- `id` (String) The unique identifier for the resource (set to the session ID).
- `created` (String) Timestamp when the session was created (RFC3339).
- `import_type` (String) The type of import. All imports are `INCREMENTAL`.
- `last_updated` (String) Timestamp when the session was last updated (RFC3339).
- `status` (String) The current status of the identity source session after the import is triggered.

## Import

Import using `{identity_source_id}/{session_id}`:

```shell
terraform import okta_identity_source_session_import.example <identity-source-id>/<session-id>
```
