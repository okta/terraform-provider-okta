---
layout: "okta"
page_title: "Okta: okta_app_metadata_saml"
sidebar_current: "docs-okta-datasource-app-metadata-saml"
description: |-
  Get a SAML application's metadata from Okta.
---

# okta_app_metadata_saml

Use this data source to retrieve the collaborators for a given repository.

## Example Usage

```hcl
data "okta_app_metadata_saml" "example" {
  app_id = "<app id>"
  key_id = "<cert key id>"
}
```

## Arguments Reference

 * `app_id` - (Required) The application ID.

 * `key_id` - (Required) Certificate Key ID.

## Attributes Reference

 * `metadata` - raw metadata of application.

 * `http_redirect_binding` - urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect location from the SAML metadata.

 * `http_post_binding` - urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Post location from the SAML metadata.

 * `certificate` - public certificate from application metadata.

 * `want_authn_requests_signed` - Whether authn requests are signed.

 * `entity_id` - Entity URL for instance `https://www.okta.com/saml2/service-provider/sposcfdmlybtwkdcgtuf`.
