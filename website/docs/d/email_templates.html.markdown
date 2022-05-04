---
layout: 'okta'
page_title: 'Okta: okta_email_templates'
sidebar_current: 'docs-okta-datasource-email-templates'
description: |-
Get the email templates belonging to a brand in an Okta organization.
---


# okta_email_templates

Use this data source to retrieve the [email
templates](https://developer.okta.com/docs/reference/api/brands/#email-template)
of a brand in an Okta organization.

## Example Usage

```hcl
data "okta_brands" "test" {
}

data "okta_email_templates" "test" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
}
```

## Argument Reference

- `brand_id` - (Required) Brand ID

## Attribute Reference

- `email_templates` - List of `okta_email_template` belonging to the brand
