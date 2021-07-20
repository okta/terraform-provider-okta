---
layout: 'okta'
page_title: 'Okta: okta_idp_saml_key'
sidebar_current: 'docs-okta-resource-saml-key'
description: |-
  Creates a SAML Identity Provider Signing Key.
---

# okta_idp_saml_key

Creates a SAML Identity Provider Signing Key.

This resource allows you to create and configure a SAML Identity Provider Signing Key.

## Example Usage

```hcl
resource "okta_idp_saml_key" "example_1" {
  x5c = ["${okta_app_saml.example.certificate}"]
}

resource "okta_idp_saml_key" "example_2" {
  x5c = ["MIIDqDCCApCgAwIBAgIGAVGNO4qeMA0GCSqGSIb3DQEBBQUAMIGUMQswCQYDVQQGEwJVUzETMBEG\nA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEU\nMBIGA1UECwwLU1NPUHJvdmlkZXIxFTATBgNVBAMMDGJhbGFjb21wdGVzdDEcMBoGCSqGSIb3DQEJ\nARYNaW5mb0Bva3RhLmNvbTAeFw0xNTEyMTAxODUwMDhaFw0xNzEyMTAxODUxMDdaMIGUMQswCQYD\nVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsG\nA1UECgwET2t0YTEUMBIGA1UECwwLU1NPUHJvdmlkZXIxFTATBgNVBAMMDGJhbGFjb21wdGVzdDEc\nMBoGCSqGSIb3DQEJARYNaW5mb0Bva3RhLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC\nggEBALAakG48bgcTWHdwmVLHig0mkiRejxIVm3wbzrNSJcBruTq2zCYZ1rGfVxTYON8kJqvkXPmv\nkzWKhpEkvhubL+mx29XpXY0AsNIfgcm5xIV56yhXSvlMdqzGo3ciRwoACaF+ClNLxmXK9UTZD89B\nbVVGCG5AEvja0eCQ0GYsO5i9aSI5aTroab8Aew31PuWl/RGQWmjVy8+7P4wwkKKJNKCpxMYDlhfa\nWRp0zwUSbUCO0qEyeAYdZx6CLES4FGrDi/7D6G+ewWC+kbz1tL1XpF2Dcg3+IOlHrV6VWzz3rG39\nv9zFIncjvoQJFDGWhpqGqcmXvgH0Ze3SVcVF01T+bK0CAwEAATANBgkqhkiG9w0BAQUFAAOCAQEA\nAHmnSZ4imjNrIf9wxfQIcqHXEBoJ+oJtd59cw1Ur/YQY9pKXxoglqCQ54ZmlIf4GghlcZhslLO+m\nNdkQVwSmWMh6KLxVM18/xAkq8zyKbMbvQnTjFB7x45bgokwbjhivWqrB5LYHHCVN7k/8mKlS4eCK\nCi6RGEmErjojr4QN2xV0qAqP6CcGANgpepsQJCzlWucMFKAh0x9Kl8fmiQodfyLXyrebYsVnLrMf\njxE1b6dg4jKvv975tf5wreQSYZ7m//g3/+NnuDKkN/03HqhV7hTNi1fyctXk8I5Nwgyr+pT5LT2k\nYoEdncuy+GQGzE9yLOhC4HNfHQXpqp2tMPdRlw=="]
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
$ terraform import okta_idp_saml_key.example <key id>
```
