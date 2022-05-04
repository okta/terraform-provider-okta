---
layout: 'okta'
page_title: 'Okta: okta_email_customization'
sidebar_current: 'docs-okta-datasource-email-customization'
description: |-
Get the email customization of an email template belonging to a brand in an Okta organization.
---

# okta_email_customization

Use this data source to retrieve the [email
customization](https://developer.okta.com/docs/reference/api/brands/#get-email-customization)
of an email template belonging to a brand in an Okta organization.

## Example Usage

```hcl
data "okta_brands" "test" {
}

data "okta_email_customizations" "forgot_password" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
}

data "okta_email_customization" "forgot_password_en" {
  customization_id = tolist(data.okta_email_customizations.forgot_password.email_customizations)[0].id
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
}
```

## Arguments Reference

- `customization_id` - (Required) Customization ID
- `brand_id` - (Required) Brand ID
- `template_name` - (Required) Template Name

## Attributes Reference

- `id` - (Required) Customization ID
- `links` - Link relations for this object - JSON HAL - Discoverable resources related to the email template
- `language` - The language supported by the customization
- `is_default` - Whether the customization is the default
- `subject` - The subject of the customization
- `body` - The body of the customization
