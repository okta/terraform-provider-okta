---
page_title: "okta_identity_source_bulk_upsert Resource - terraform-provider-okta"
description: |-
  Manages an Okta Identity Source Bulk Upsert resource.
---

# okta_identity_source_bulk_upsert (Resource)

Manages an Okta Identity Source Bulk Upsert.

## Example Usage

```hcl
resource "okta_identity_source_bulk_upsert" "example" {
  identity_source_id = "<identity-source-id>"
  session_id = "<session-id>"

  # Optional fields
  # entity_type = "<entity_type>"
  # external_id = "<external-id>"
  # email = "user@example.com"
  # first_name = "Example First Name"
  # home_address = "<home_address>"
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source
- `session_id` (String) ID of the identity source session

### Optional

- `entity_type` (String) The type of data to upsert into the session. Currently, only `USERS` is supported.
- `external_id` (String) The external ID of the entity that needs to be created or updated in Okta
- `email` (String) Email address of the user
- `first_name` (String) First name of the user
- `home_address` (String) Home address of the user
- `last_name` (String) Last name of the user
- `mobile_phone` (String) Mobile phone number of the user
- `second_email` (String) Alternative email address of the user
- `user_name` (String) Username of the user

### Read-Only

- `id` (String) The unique identifier for the resource.

## Import

Import using `{identity_source_id}/{session_id}/{id}`:

```shell
terraform import okta_identity_source_bulk_upsert.example <identity_source_id> <session_id> <id>
```
