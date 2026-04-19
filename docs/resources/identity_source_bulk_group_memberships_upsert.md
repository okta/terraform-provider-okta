---
page_title: "Resource: okta_identity_source_bulk_group_memberships_upsert"
description: |-
  Uploads a batch of group membership assignments to an Okta Identity Source session for creation or update.
---

# Resource: okta_identity_source_bulk_group_memberships_upsert

Uploads a batch of group membership assignments into an Okta Identity Source session, staging them for creation or update in Okta. Use this resource alongside `okta_identity_source_session` and `okta_identity_source_session_import` to associate users with groups. The operation is write-only: Okta does not expose a GET endpoint for the uploaded data, so state is preserved from the create call.

~> **Note:** This resource does not support deletion. Removing it from configuration takes no action against Okta.

## Example Usage

```terraform
# Session 1: upload user and group, then import to commit them to Okta
resource "okta_identity_source_session" "session_setup" {
  identity_source_id = "<identity-source-id>"
}

resource "okta_identity_source_bulk_upsert" "user_setup" {
  identity_source_id = okta_identity_source_session.session_setup.identity_source_id
  session_id         = okta_identity_source_session.session_setup.id
  entity_type        = "USERS"

  profiles {
    external_id = "USEREXT123456EXAMPLE"

    profile {
      user_name  = "jdoe@example.com"
      email      = "jdoe@example.com"
      first_name = "Jane"
      last_name  = "Doe"
    }
  }
}

resource "okta_identity_source_group" "group_setup" {
  identity_source_id = "<identity-source-id>"
  external_id        = "GROUPEXT123456EXAMPLE"

  profile {
    display_name = "Engineering"
    description  = "Engineering team"
  }
}

resource "okta_identity_source_session_import" "session_import_setup" {
  identity_source_id = okta_identity_source_session.session_setup.identity_source_id
  session_id         = okta_identity_source_session.session_setup.id

  depends_on = [okta_identity_source_bulk_upsert.user_setup]
}

# Session 2: upload group memberships and import
resource "okta_identity_source_session" "example" {
  identity_source_id = "<identity-source-id>"
  depends_on         = [okta_identity_source_session_import.session_import_setup]
}

resource "okta_identity_source_bulk_group_memberships_upsert" "example" {
  identity_source_id = okta_identity_source_session.example.identity_source_id
  session_id         = okta_identity_source_session.example.id

  memberships {
    group_external_id   = okta_identity_source_group.group_setup.external_id
    member_external_ids = ["USEREXT123456EXAMPLE"]
  }
}

resource "okta_identity_source_session_import" "example" {
  identity_source_id = okta_identity_source_session.example.identity_source_id
  session_id         = okta_identity_source_session.example.id

  depends_on = [okta_identity_source_bulk_group_memberships_upsert.example]
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source. Forces replacement when changed.
- `session_id` (String) ID of the identity source session to upload data into. Forces replacement when changed.

### Optional

- `memberships` (Block List, Optional) Array of group memberships to insert or update in Okta. (see [below for nested schema](#nested-schema-for-memberships))

### Read-Only

- `id` (String) The unique identifier for the resource (set to the session ID).

### Nested Schema for `memberships`

Optional:

- `group_external_id` (String) The external ID of the group whose memberships need to be inserted or updated in Okta.
- `member_external_ids` (List of String) Array of external IDs of member profiles to insert into the group in Okta.

## Import

Import using `{identity_source_id}/{session_id}/{id}`:

```shell
terraform import okta_identity_source_bulk_group_memberships_upsert.example <identity-source-id>/<session-id>/<id>
```
