---
page_title: "okta_brand_themes Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Brand Themes.
---

# okta_brand_themes (Data Source)

Use this data source to retrieve an Okta Brand Themes.

## Example Usage

```hcl
data "okta_brand_themes" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
