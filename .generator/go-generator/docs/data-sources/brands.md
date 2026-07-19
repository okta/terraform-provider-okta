---
page_title: "okta_brands Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Brands.
---

# okta_brands (Data Source)

Use this data source to retrieve an Okta Brands.

## Example Usage

```hcl
data "okta_brands" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
