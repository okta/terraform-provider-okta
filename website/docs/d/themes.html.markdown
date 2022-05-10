---
layout: 'okta'
page_title: 'Okta: okta_themes'
sidebar_current: 'docs-okta-datasource-themes'
description: |-
Get Themes of a Brand of an Okta Organization.
---


# okta_themes

Use this data source to retrieve 
[Themes](https://developer.okta.com/docs/reference/api/brands/#theme-response-object)
of a brand for an Okta orgnanization.

## Example Usage

```hcl
data "okta_brands" "test" {
}

data "okta_themes" "example" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
}
```

## Argument Reference

- `brand_id` - (Required) Brand ID

## Attribute Reference

- `themes` - List of `okta_theme` belonging to the brand.
