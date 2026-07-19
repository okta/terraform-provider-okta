---
page_title: "okta_application_token Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Token resource.
---

# okta_application_token (Resource)

Manages an Okta Application Token.

## Example Usage

```hcl
resource "okta_application_token" "example" {
  app_id = "<app-id>"

  # Optional fields
  # client_id = "<client-id>"
  # issuer = "<issuer>"
  # scopes = "<scopes>"
  # user_id = "<user-id>"
}
```

## Schema

### Required

- `app_id` (String) ID of the parent application

### Optional

- `client_id` (String) Client ID
- `issuer` (String) The complete URL of the authorization server that issued the Token
- `scopes` (List) The scope names attached to the Token
- `user_id` (String) The ID of the user associated with the Token

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) The embedded resources related to the object if the `expand` query parameter is specified
- `created` (String) Timestamp when the object was created
- `expires_at` (String) Expiration time of the OAuth 2.0 Token
- `last_updated` (String) Timestamp when the object was last updated
- `status` (String) Status

## Import

Import using `{app_id}/{id}`:

```shell
terraform import okta_application_token.example <app_id> <id>
```
