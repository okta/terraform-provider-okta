---
page_title: "okta_principal_rate_limit Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Principal Rate Limit.
---

# okta_principal_rate_limit (Data Source)

Use this data source to retrieve an Okta Principal Rate Limit.

## Example Usage

```hcl
data "okta_principal_rate_limit" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
