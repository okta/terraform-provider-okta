---
page_title: "Resource: okta_identity_source_session"
description: |-
  Manages an Okta Identity Source Session.
---

# Resource: okta_identity_source_session

Manages an Okta Identity Source Session. A session represents a single import operation against an identity source and must be created before uploading user or group data.

## Example Usage

```terraform
resource "okta_identity_source_session" "example" {
  identity_source_id = "<identity-source-id>"
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source. Forces replacement when changed.

### Read-Only

- `id` (String) The unique identifier of the identity source session.
- `created` (String) Timestamp when the session was created (RFC3339).
- `import_type` (String) The type of import. All imports are `INCREMENTAL`.
- `last_updated` (String) Timestamp when the session was last updated (RFC3339).
- `status` (String) The current status of the identity source session.

## Import

Import using `{identity_source_id}/{session_id}`:

```shell
terraform import okta_identity_source_session.example <identity-source-id>/<session-id>
```
