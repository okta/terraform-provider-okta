---
page_title: "okta_user_factors Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta User Factors.
---

# okta_user_factors (Data Source)

Use this data source to retrieve an Okta User Factors.

## Example Usage

```hcl
data "okta_user_factors" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
