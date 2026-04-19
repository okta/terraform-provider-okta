---
page_title: "Resource: okta_identity_source_bulk_groups_upsert"
description: |-
  Uploads a batch of group profiles to an Okta Identity Source session for creation or update.
---

# Resource: okta_identity_source_bulk_groups_upsert

Uploads a batch of group profiles into an Okta Identity Source session, staging them for creation or update in Okta. Use this resource alongside `okta_identity_source_session`. The operation is write-only: Okta does not expose a GET endpoint for the uploaded data, so state is preserved from the create call.

~> **Note:** This resource does not support deletion. Removing it from configuration takes no action against Okta.

## Example Usage

```terraform
resource "okta_identity_source_session" "example" {
  identity_source_id = "<identity-source-id>"
}

resource "okta_identity_source_bulk_groups_upsert" "example" {
  identity_source_id = okta_identity_source_session.example.identity_source_id
  session_id         = okta_identity_source_session.example.id

  profiles {
    external_id = "GROUPEXT123456EXAMPLE"

    group_profile {
      display_name = "Engineering"
      description  = "Engineering team group"
    }
  }
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source. Forces replacement when changed.
- `session_id` (String) ID of the identity source session to upload data into. Forces replacement when changed.

### Optional

- `profiles` (Block List, Optional) Array of group profiles to create or update in Okta. (see [below for nested schema](#nested-schema-for-profiles))

### Read-Only

- `id` (String) The unique identifier for the resource (set to the session ID).

### Nested Schema for `profiles`

Optional:

- `external_id` (String) The external ID of the group to create or update in Okta.
- `group_profile` (Block, Optional) Group profile attributes. (see [below for nested schema](#nested-schema-for-profilesgroup_profile))

### Nested Schema for `profiles.group_profile`

Optional:

- `display_name` (String) Name of the group.
- `description` (String) Description of the group.

## Import

Import using `{identity_source_id}/{session_id}/{id}`:

```shell
terraform import okta_identity_source_bulk_groups_upsert.example <identity-source-id>/<session-id>/<id>
```
