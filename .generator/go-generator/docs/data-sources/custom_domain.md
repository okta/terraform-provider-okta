---
page_title: "okta_custom_domain Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Custom Domain.
---

# okta_custom_domain (Data Source)

Use this data source to retrieve an Okta Custom Domain.

## Example Usage

```hcl
data "okta_custom_domain" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
