---
layout: 'okta'
page_title: 'Okta: okta_idp_saml_key'
sidebar_current: 'docs-okta-resource-saml-key'
description: |-
  Creates a SAML Identity Provider Signing Key.
---

# okta_idp_saml_key

This resource allows you to create and configure a SAML Identity Provider Signing Key.

## IMPORTANT NOTE

Identity Provider Signing Key can not be updated, it can only be created or removed. Thus, in situation
where Identity Provider Signing Key should be updated, one can not simply change `x5c` and apply the changes.
This will cause the error, because the key can not be removed if it's used by the provider. To avoid this situation,
the behavior of the resource update process is as follows:

- create a new identity provider key
- list all identity providers that are using the old key (**even the ones that are managed outside terraform!**)
- assign a new key to these identity providers
- remove old identity provider key

For more details, please refer to the [original issue](https://github.com/okta/terraform-provider-okta/issues/672). 

## Example Usage

```hcl
resource "okta_idp_saml_key" "example_1" {
  x5c = ["${okta_app_saml.example.certificate}"]
}

resource "okta_idp_saml_key" "example_2" {
  x5c = ["MIIDnjCCAoagAwIBAgIGAVG3MN+PMA0GCSqGSIb3DQEBBQUAMIGPMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5p\nYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEUMBIGA1UECwwLU1NPUHJvdmlkZXIxEDAOBgNVBAMM\nB2V4YW1wbGUxHDAaBgkqhkiG9w0BCQEWDWluZm9Ab2t0YS5jb20wHhcNMTUxMjE4MjIyMjMyWhcNMjUxMjE4MjIyMzMyWjCB\njzELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNVBAoMBE9r\ndGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRAwDgYDVQQDDAdleGFtcGxlMRwwGgYJKoZIhvcNAQkBFg1pbmZvQG9rdGEuY29t\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtcnyvuVCrsFEKCwHDenS3Ocjed8eWDv3zLtD2K/iZfE8BMj2wpTf\nn6Ry8zCYey3mWlKdxIybnV9amrujGRnE0ab6Q16v9D6RlFQLOG6dwqoRKuZy33Uyg8PGdEudZjGbWuKCqqXEp+UKALJHV+k4\nwWeVH8g5d1n3KyR2TVajVJpCrPhLFmq1Il4G/IUnPe4MvjXqB6CpKkog1+ThWsItPRJPAM+RweFHXq7KfChXsYE7Mmfuly8s\nDQlvBmQyxZnFHVuiPfCvGHJjpvHy11YlHdOjfgqHRvZbmo30+y0X/oY/yV4YEJ00LL6eJWU4wi7ViY3HP6/VCdRjHoRdr5L/\nDwIDAQABMA0GCSqGSIb3DQEBBQUAA4IBAQCzzhOFkvyYLNFj2WDcq1YqD4sBy1iCia9QpRH3rjQvMKDwQDYWbi6EdOX0TQ/I\nYR7UWGj+2pXd6v0t33lYtoKocp/4lUvT3tfBnWZ5KnObi+J2uY2teUqoYkASN7F+GRPVOuMVoVgm05ss8tuMb2dLc9vsx93s\nDt+XlMTv/2qi5VPwaDtqduKkzwW9lUfn4xIMkTiVvCpe0X2HneD2Bpuao3/U8Rk0uiPfq6TooWaoW3kjsmErhEAs9bA7xuqo\n1KKY9CdHcFhkSsMhoeaZylZHtzbnoipUlQKSLMdJQiiYZQ0bYL83/Ta9fulr1EERICMFt3GUmtYaZZKHpWSfdJp9"]
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
