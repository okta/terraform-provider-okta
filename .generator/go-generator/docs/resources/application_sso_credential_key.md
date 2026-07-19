---
page_title: "okta_application_sso_credential_key Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Sso Credential Key resource.
---

# okta_application_sso_credential_key (Resource)

Manages an Okta Application Sso Credential Key.

## Example Usage

```hcl
resource "okta_application_sso_credential_key" "example" {
  app_id = "<app-id>"

  # Optional fields
  # n = "<n>"
}
```

## Schema

### Required

- `app_id` (String) ID of the parent application

### Optional

- `n` (String) RSA modulus value that is used by both the public and private keys and provides a link between them

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the object was created
- `e` (String) RSA key value (public exponent) for Key binding
- `expires_at` (String) Timestamp when the certificate expires
- `kid` (String) Unique identifier for the certificate
- `kty` (String) Cryptographic algorithm family for the certificate's keypair. Valid value: `RSA`
- `last_updated` (String) Timestamp when the object was last updated
- `use` (String) Acceptable use of the certificate. Valid value: `sig`
- `x5c` (List) X.509 certificate chain that contains a chain of one or more certificates
- `x5t_s256` (String) X.509 certificate SHA-256 thumbprint, which is the base64url-encoded SHA-256 thumbprint (digest) of the DER encoding of an X.509 certificate

## Import

Import using `{app_id}/{id}`:

```shell
terraform import okta_application_sso_credential_key.example <app_id> <id>
```
