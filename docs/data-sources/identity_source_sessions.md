---
page_title: "Data Source: okta_identity_source_sessions"
description: |-
  Retrieves an Okta Identity Source Session.
---

# Data Source: okta_identity_source_sessions

Retrieves an identity source session by ID, or the most recently created session for the given identity source if no `id` is specified.

## Example Usage

```terraform
# Look up the most recent session
data "okta_identity_source_sessions" "latest" {
  identity_source_id = "<identity-source-id>"
}

# Look up a specific session by ID
data "okta_identity_source_sessions" "by_id" {
  identity_source_id = "<identity-source-id>"
  id                 = "<session-id>"
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source.

### Optional

- `id` (String) Unique identifier of the identity source session to look up. If omitted, the most recently created session is returned.

### Read-Only

- `created` (String) Timestamp when the session was created (RFC3339).
- `import_type` (String) The type of import. All imports are `INCREMENTAL`.
- `last_updated` (String) Timestamp when the session was last updated (RFC3339).
- `status` (String) The current status of the identity source session.
