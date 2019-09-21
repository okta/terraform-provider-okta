---
layout: "okta"
page_title: "Okta: okta_idp_saml"
sidebar_current: "docs-okta-datasource-idp-saml"
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

 * `label` - (Optional) The label of the app to retrieve, conflicts with `label_prefix` and `id`.

 * `label_prefix` - (Optional) Label prefix of the app to retrieve, conflicts with `label` and `id`. This will tell the provider to do a `starts with` query as opposed to an `equals` query.

 * `id` - (Optional) `id` of application to retrieve, conflicts with `label` and `label_prefix`.

 * `active_only` - (Optional) tells the provider to query for only `ACTIVE` applications.

## Attributes Reference

 * `id` - `id` of application.

 * `label` - `label` of application.

 * `description` - `description` of application.

 * `name` - `name` of application.

 * `status` - `status` of application.
