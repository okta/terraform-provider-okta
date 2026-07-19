---
page_title: "okta_identity_source_user Resource - terraform-provider-okta"
description: |-
  Manages an Okta Identity Source User resource.
---

# okta_identity_source_user (Resource)

Manages an Okta Identity Source User.

## Example Usage

```hcl
resource "okta_identity_source_user" "example" {
  identity_source_id = "<identity-source-id>"

  # Optional fields
  # email = "user@example.com"
  # first_name = "Example First Name"
  # home_address = "<home_address>"
  # last_name = "Example Last Name"
  # mobile_phone = "<mobile_phone>"
}
```

## Schema

### Required

- `identity_source_id` (String) ID of the identity source

### Optional

- `email` (String) Email address of the user
- `first_name` (String) First name of the user
- `home_address` (String) Home address of the user
- `last_name` (String) Last name of the user
- `mobile_phone` (String) Mobile phone number of the user
- `second_email` (String) Alternative email address of the user
- `user_name` (String) Username of the user

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) The timestamp when the user was created in the identity source
- `external_id` (String) The external ID of the user in the identity source
- `last_updated` (String) The timestamp when the user was last updated in the identity source

## Import

Import using `{identity_source_id}/{id}`:

```shell
terraform import okta_identity_source_user.example <identity_source_id> <id>
```
