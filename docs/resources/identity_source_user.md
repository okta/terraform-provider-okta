---
page_title: "Resource: okta_identity_source_user"
description: |-
  Manages an Okta Identity Source User.
---

# Resource: okta_identity_source_user

Manages an individual user in an Okta Identity Source. This resource creates or updates a single user record directly in the identity source (without requiring a session), using the external ID as the unique key.

## Example Usage

```terraform
resource "okta_identity_source_user" "example" {
  identity_source_id = "<identity-source-id>"
  id                 = "USEREXT123456EXAMPLE"

  profile {
    user_name  = "jdoe@example.com"
    email      = "jdoe@example.com"
    first_name = "Jane"
    last_name  = "Doe"
  }
}
```

## Schema

### Required

- `id` (String) The external ID of the user in the identity source. Used as the resource identifier. Forces replacement when changed.
- `identity_source_id` (String) ID of the identity source. Forces replacement when changed.

### Optional

- `profile` (Block, Optional) User profile attributes. (see [below for nested schema](#nested-schema-for-profile))

### Read-Only

- `created` (String) Timestamp when the user was created in the identity source (RFC3339).
- `last_updated` (String) Timestamp when the user was last updated in the identity source (RFC3339).

### Nested Schema for `profile`

Optional:

- `email` (String) Email address of the user.
- `first_name` (String) First name of the user.
- `home_address` (String) Home address of the user.
- `last_name` (String) Last name of the user.
- `mobile_phone` (String) Mobile phone number of the user.
- `second_email` (String) Alternative email address of the user.
- `user_name` (String) Username of the user.

## Import

Import using `{identity_source_id}/{id}`:

```shell
terraform import okta_identity_source_user.example <identity-source-id>/<external-id>
```