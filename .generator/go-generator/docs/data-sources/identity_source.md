---
page_title: "okta_identity_source Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Identity Source.
---

# okta_identity_source (Data Source)

Use this data source to retrieve an Okta Identity Source.

## Example Usage

```hcl
data "okta_identity_source" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
