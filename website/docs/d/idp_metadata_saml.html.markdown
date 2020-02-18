---
layout: "okta"
page_title: "Okta: okta_idp_metadata_saml"
sidebar_current: "docs-okta-datasource-idp-metadata-saml"
description: |-
  Get SAML IdP metadata from Okta.
---

# okta_idp_metadata_saml

Use this data source to retrieve SAML IdP metadata from Okta.

## Example Usage

```hcl
data "okta_idp_metadata_saml" "example" {
  id = "<idp id>"
}
```

## Arguments Reference

* `idp_id` - (Required) The id of the IdP to retrieve metadata for.

## Attributes Reference

* `assertions_signed` - whether assertions are signed.

* `authn_request_signed` - whether authn requests are signed.

* `encryption_certificate` - SAML request encryption certificate.

* `entity_id` - Entity URL for instance `https://www.okta.com/saml2/service-provider/sposcfdmlybtwkdcgtuf`.

* `http_post_binding` - urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Post location from the SAML metadata.

* `http_redirect_binding` - urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect location from the SAML metadata.

* `metadata` - raw IdP metadata.

* `signing_certificate` - SAML request signing certificate.
