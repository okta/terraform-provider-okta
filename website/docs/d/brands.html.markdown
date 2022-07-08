---
layout: 'okta'
page_title: 'Okta: okta_brands'
sidebar_current: 'docs-okta-datasource-brands'
description: |-
Get the brands belonging to an Okta organization.
---


# okta_brands

Use this data source to retrieve the brands belonging to an Okta organization.

## Example Usage

```hcl
data "okta_brands" "test" {
}
```

## Argument Reference

No arguments.

## Attribute Reference

- `brands` - List of `okta_brand` belonging to the organization
