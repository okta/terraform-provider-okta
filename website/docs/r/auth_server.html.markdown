---
layout: 'okta'
page_title: 'Okta: okta_auth_server'
sidebar_current: 'docs-okta-resource-auth-server'
description: |-
  Creates an Authorization Server.
---

# okta_auth_server

Creates an Authorization Server.

This resource allows you to create and configure an Authorization Server.

## Example Usage

```hcl
resource "okta_auth_server" "example" {
  audiences   = ["api://example"]
  description = "My Example Auth Server"
  name        = "example"
  issuer_mode = "CUSTOM_URL"
  status      = "ACTIVE"
}
```

## Argument Reference

The following arguments are supported:

- `audiences` - (Required) The recipients that the tokens are intended for. This becomes the `aud` claim in an access token.

- `status` - (Optional) The status of the auth server. It defaults to `"ACTIVE"`

- `credentials_rotation_mode` - (Optional) The key rotation mode for the authorization server. Can be `"AUTO"` or `"MANUAL"`.

- `description` - (Optional) The description of the authorization server.

- `name` - (Required) The name of the authorization server.

- `issuer_mode` - (Optional) Allows you to use a custom issuer URL. It can be set to `"CUSTOM_URL"` or `"ORG_URL"`

## Attributes Reference

- `id` - ID of the authorization server.

- `kid` - The ID of the JSON Web Key used for signing tokens issued by the authorization server.

- `issuer` - The complete URL for a Custom Authorization Server. This becomes the `iss` claim in an access token.

- `credentials_last_rotated` - The timestamp when the authorization server started to use the `kid` for signing tokens.

- `credentials_next_rotation` - The timestamp when the authorization server changes the key for signing tokens. Only returned when `credentials_rotation_mode` is `"AUTO"`.

## Import

Authorization Server can be imported via the Okta ID.

```
$ terraform import okta_auth_server.example <auth server id>
```
