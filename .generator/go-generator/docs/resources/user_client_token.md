---
page_title: "okta_user_client_token Resource - terraform-provider-okta"
description: |-
  Manages an Okta User Client Token resource.
---

# okta_user_client_token (Resource)

Manages an Okta User Client Token.

## Example Usage

```hcl
resource "okta_user_client_token" "example" {
  user_id = "<user-id>"
  client_id = "<client-id>"

  # Optional fields
  # issuer = "<issuer>"
  # scopes = "<scopes>"
}
```

## Schema

### Required

- `user_id` (String) ID of the user
- `client_id` (String) ID of the OAuth client

### Optional

- `issuer` (String) The complete URL of the authorization server that issued the Token
- `scopes` (List) The scope names attached to the Token

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) The embedded resources related to the object if the `expand` query parameter is specified
- `created` (String) Timestamp when the object was created
- `expires_at` (String) Expiration time of the OAuth 2.0 Token
- `last_updated` (String) Timestamp when the object was last updated
- `status` (String) Status

## Import

Import using `{user_id}/{client_id}/{id}`:

```shell
terraform import okta_user_client_token.example <user_id> <client_id> <id>
```
