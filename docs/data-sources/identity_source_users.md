---
page_title: "Data Source: okta_identity_source_users"
description: |-
  Retrieves an Okta Identity Source User.
---

# Data Source: okta_identity_source_users

Retrieves a user record from an Okta Identity Source by external ID.

## Example Usage

```terraform
data "okta_identity_source_users" "example" {
  identity_source_id = "<identity-source-id>"
  external_id        = "USEREXT123456EXAMPLE"
}

output "user_email" {
  value = data.okta_identity_source_users.example.profile.email
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source.
- `external_id` (String) The external ID of the user in the identity source.

### Optional

- `id` (String) Unique identifier of the user. If provided, takes precedence over `external_id` for the lookup.

### Read-Only

- `created` (String) Timestamp when the user was created in the identity source (RFC3339).
- `last_updated` (String) Timestamp when the user was last updated in the identity source (RFC3339).
- `profile` (Object) User profile attributes.
  - `email` (String) Email address of the user.
  - `first_name` (String) First name of the user.
  - `home_address` (String) Home address of the user.
  - `last_name` (String) Last name of the user.
  - `mobile_phone` (String) Mobile phone number of the user.
  - `second_email` (String) Alternative email address of the user.
  - `user_name` (String) Username of the user.
