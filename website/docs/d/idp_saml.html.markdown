---
layout: 'okta'
page_title: 'Okta: okta_idp_saml'
sidebar_current: 'docs-okta-datasource-idp-saml'
description: |-
  Get a SAML IdP from Okta.
---

# okta_idp_saml

Use this data source to retrieve a SAML IdP from Okta.

## Example Usage

```hcl
data "okta_idp_saml" "example" {
  label = "Example App"
}
```

## Arguments Reference

- `name` - (Optional) The name of the idp to retrieve, conflicts with `id`.

- `id` - (Optional) The id of the idp to retrieve, conflicts with `name`.

## Attributes Reference

- `id` - id of idp.

- `name` - name of the idp.

- `type` - type of idp.

- `acs_binding` - HTTP binding used to receive a SAMLResponse message from the IdP.

- `acs_type` - Determines whether to publish an instance-specific (trust) or organization (shared) ACS endpoint in the SAML metadata.

- `sso_url` - single sign on url.

- `sso_binding` - single sign on binding.

- `sso_destination` - SSO request binding, HTTP-POST or HTTP-REDIRECT.

- `subject_format` - Expression to generate or transform a unique username for the IdP user.

- `subject_filter` - regular expression pattern used to filter untrusted IdP usernames.

- `issuer` - URI that identifies the issuer (IdP).

- `issuer_mode` - indicates whether Okta uses the original Okta org domain URL, or a custom domain URL in the request to the IdP.

- `audience` - URI that identifies the target Okta IdP instance (SP)

- `kid` - Key ID reference to the IdP's X.509 signature certificate.
