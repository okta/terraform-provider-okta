---
layout: 'okta'
page_title: 'Okta: okta_idp_saml_signing_key'
sidebar_current: 'docs-okta-resource-saml-signing-key'
description: |-
  Creates a SAML Identity Provider Signing Key.
---

# okta_idp_saml_signing_key

Creates a SAML Identity Provider Signing Key.

This resource allows you to create and configure a SAML Identity Provider Signing Key.

## Example Usage

```hcl
resource "okta_idp_saml_key" "example" {
  x5c = ["${okta_app_saml.example.certificate}"]
}
```

## Argument Reference

The following arguments are supported:

- `x5c` - (Required) base64-encoded X.509 certificate chain with DER encoding.

## Attributes Reference

- `id` - Key ID.

- `kid` - Key ID.

- `created` - Date created.

- `expires_at` - Date the cert expires.

- `kty` - Identifies the cryptographic algorithm family used with the key.

- `use` - Intended use of the public key.

- `x5t_s256` - base64url-encoded SHA-256 thumbprint of the DER encoding of an X.509 certificate.

## Import

A SAML IdP Signing Key can be imported via the key id.

```
$ terraform import okta_idp_saml_signing_key.example <key id>
```
