---
layout: 'okta'
page_title: 'Okta: okta_brand'
sidebar_current: 'docs-okta-datasource-brand'
description: |-
  Get a single Brand from Okta.
---

# okta_brand

Use this data source to retrieve a [Brand](https://developer.okta.com/docs/reference/api/brands/#brand-object) from Okta.

## Example Usage

```hcl
data "okta_brands" "test" {
}

data "okta_brand" "test" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
}
```

## Arguments Reference

- `brand_id` - (Required) Brand ID

## Attributes Reference

- `id` - Brand ID

- `custom_privacy_policy_url` - Custom privacy policy URL

- `links` - Link relations for this object - JSON HAL - Discoverable resources related to the brand

- `remove_powered_by_okta` - Removes "Powered by Okta" from the Okta-hosted sign-in page, and "© 2021 Okta, Inc." from the Okta End-User Dashboard
