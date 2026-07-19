---
page_title: "okta_identity_provider_signing_key Resource - terraform-provider-okta"
description: |-
  Manages an Okta Identity Provider Signing Key resource.
---

# okta_identity_provider_signing_key (Resource)

Manages an Okta Identity Provider Signing Key.

## Example Usage

```hcl
resource "okta_identity_provider_signing_key" "example" {
  idp_id = "<idp-id>"

  # Optional fields
  # e = "<e>"
  # kid = "<kid>"
  # kty = "<kty>"
  # n = "<n>"
  # use = "<use>"
}
```

## Schema

### Required

- `idp_id` (String) ID of the identity provider

### Optional

- `e` (String) The exponent value for the RSA public key
- `kid` (String) Unique identifier for the key
- `kty` (String) Identifies the cryptographic algorithm family used with the key
- `n` (String) The modulus value for the RSA public key
- `use` (String) Intended use of the public key
- `x5c` (List) Base64-encoded X.509 certificate chain with DER encoding
- `x5t_s256` (String) Base64url-encoded SHA-256 thumbprint of the DER encoding of an X.509 certificate

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the object was created
- `expires_at` (String) Timestamp when the object expires
- `last_updated` (String) Timestamp when the object was last updated

## Import

Import using `{idp_id}/{id}`:

```shell
terraform import okta_identity_provider_signing_key.example <idp_id> <id>
```
