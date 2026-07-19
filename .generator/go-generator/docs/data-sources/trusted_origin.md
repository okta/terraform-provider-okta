---
page_title: "okta_trusted_origin Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Trusted Origin.
---

# okta_trusted_origin (Data Source)

Use this data source to retrieve an Okta Trusted Origin.

## Example Usage

```hcl
data "okta_trusted_origin" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
