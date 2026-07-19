---
page_title: "okta_authorization_server Resource - terraform-provider-okta"
description: |-
  Manages an Okta Authorization Server resource.
---

# okta_authorization_server (Resource)

Manages an Okta Authorization Server.

## Example Usage

```hcl
resource "okta_authorization_server" "example" {

  # Optional fields
  # access_token_encrypted_response_algorithm = "<access_token_encrypted_response_algorithm>"
  # audiences = "<audiences>"
  # rotation_mode = "<rotation_mode>"
  # use = "<use>"
  # description = "Example description"
}
```

## Schema

### Optional

- `access_token_encrypted_response_algorithm` (String) The algorithm for encrypting access tokens issued by this authorization server. If this is requested, the response is signed, and then encrypted. The result is a nested JWT. The default, if omitted...
- `audiences` (List) The recipients that the tokens are intended for. This becomes the `aud` claim in an access token. Okta currently supports only one audience.
- `rotation_mode` (String) The Key rotation mode for the authorization server
- `use` (String) How the key is used
- `description` (String) The description of the custom authorization server
- `issuer` (String) The complete URL for the custom authorization server. This becomes the `iss` claim in an access token.
- `issuer_mode` (String) Indicates which value is specified in the issuer of the tokens that a custom authorization server returns: the Okta org domain URL or a custom domain URL.  `issuerMode` is visible if you have a cus...
- `e` (String) The key exponent of a RSA key
- `kid` (String) The unique identifier of the key
- `kty` (String) The type of public key
- `n` (String) The modulus of the RSA key
- `status` (String) The status of the public key
- `use` (String) The intended use of the public key
- `jwks_uri` (String) URL string that references a JSON Web Key Set for encrypting JWTs minted by the custom authorization server
- `name` (String) The name of the custom authorization server
- `status` (String) Status

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Created
- `kid` (String) The ID of the JSON Web Key used for signing tokens issued by the authorization server
- `last_rotated` (String) The timestamp when the authorization server started using the `kid` for signing tokens
- `next_rotation` (String) The timestamp when the authorization server changes the Key for signing tokens. This is only returned when `rotationMode` is set to `AUTO`.
- `last_updated` (String) LastUpdated
