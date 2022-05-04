---
layout: 'okta'
page_title: 'Okta: okta_email_customizations'
sidebar_current: 'docs-okta-datasource-email-customizations'
description: |-
Get the email customizations of an email template belonging to a brand in an Okta organization.
---


# okta_email_customizations

Use this data source to retrieve the [email
customizations](https://developer.okta.com/docs/reference/api/brands/#list-email-customizations)
of an email template belonging to a brand in an Okta organization.

## Example Usage

```hcl
data "okta_brands" "test" {
}

data "okta_email_customizations" "forgot_password" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
}
```

## Argument Reference

- `brand_id` - (Required) Brand ID
- `template_name` - (Required) Name of an Email Template

## Attribute Reference

- `email_customizations` - List of `okta_email_customization` belonging to the named email template of the brand
