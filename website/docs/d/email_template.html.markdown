---
layout: 'okta'
page_title: 'Okta: okta_email_template'
sidebar_current: 'docs-okta-datasource-email-template'
description: |-
  Get a single Email Template for a Brand belonging to an Okta organization.
---

# okta_brand

Use this data source to retrieve a specific [email
template](https://developer.okta.com/docs/reference/api/brands/#email-template)
of a brand in an Okta organization.


## Example Usage

```hcl
data "okta_brands" "test" {
}

data "okta_email_template" "forgot_password" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  name = "ForgotPassword"
}
```

## Arguments Reference

- `brand_id` - (Required) Brand ID
- `name` - (Required) Template Name

## Attributes Reference

- `links` - Link relations for this object - JSON HAL - Discoverable resources related to the email template
