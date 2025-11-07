---
page_title: "Data Source: okta_app_token"
description: |-
  Retrieves a refresh token for the specified app.
---

# Data Source: okta_app_token
 
Retrieves a refresh token for the specified app.

## Example Usage

### Basic Token Information

```terraform
data "okta_app_token" "example" {
  client_id = "0oardd5r32PWsF4421d7"
  id        = "oar1godmqw4QUiX4C1d7"
}

output "token_status" {
  value = data.okta_app_token.example.status
}

output "token_user" {
  value = data.okta_app_token.example.user_id
}
```

## Argument Reference

The following arguments are required:

* `client_id` - (Required) The unique Okta ID of the application associated with this token. This is typically the `client_id` of an application.
* `id` - (Required) The unique Okta ID of the refresh token to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_id` - The unique ID of the user associated with this token.
* `status` - The current status of the token (e.g., `ACTIVE`, `REVOKED`).
* `created` - Timestamp when the token was created, in RFC3339 format.
* `expires_at` - Timestamp when the token expires, in RFC3339 format.
* `scopes` - List of scope names attached to the token.
* `issuer` - The complete URL of the authorization server that issued the token.
