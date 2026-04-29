---
page_title: "Resource: okta_identity_source_bulk_upsert"
description: |-
  Uploads a batch of user profiles to an Okta Identity Source session for creation or update.
---

# Resource: okta_identity_source_bulk_upsert

Uploads a batch of user profiles into an Okta Identity Source session. Use this resource alongside `okta_identity_source_session` to stage users for import into Okta. The operation is write-only: Okta does not expose a GET endpoint for uploaded data, so state is preserved from the create call.

~> **Note:** This resource does not support deletion. Removing it from configuration will emit a warning but take no action against Okta.

## Example Usage

```terraform
resource "okta_identity_source_session" "example" {
  identity_source_id = "<identity-source-id>"
}

resource "okta_identity_source_bulk_upsert" "example" {
  identity_source_id = okta_identity_source_session.example.identity_source_id
  session_id         = okta_identity_source_session.example.id
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
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source. Forces replacement when changed.
- `session_id` (String) ID of the identity source session to upload data into. Forces replacement when changed.

### Optional

- `entity_type` (String) The type of data to upsert into the session. Currently, only `USERS` is supported.
- `profiles` (Block List, Optional) Array of user profiles to upload. (see [below for nested schema](#nested-schema-for-profiles))

### Read-Only

- `id` (String) The unique identifier for the resource (set to the session ID).

### Nested Schema for `profiles`

Optional:

- `external_id` (String) The external ID of the entity to create or update in Okta.
- `profile` (Block, Optional) User profile attributes. (see [below for nested schema](#nested-schema-for-profilesprofile))

### Nested Schema for `profiles.profile`

Optional:

- `email` (String) Email address of the user.
- `first_name` (String) First name of the user.
- `home_address` (String) Home address of the user.
- `last_name` (String) Last name of the user.
- `mobile_phone` (String) Mobile phone number of the user.
- `second_email` (String) Alternative email address of the user.
- `user_name` (String) Username of the user.

## Import

Import using `{identity_source_id}/{session_id}/{id}`:

```shell
terraform import okta_identity_source_bulk_upsert.example <identity-source-id>/<session-id>/<id>
```
