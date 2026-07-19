---
page_title: "okta_policy Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Policy.
---

# okta_policy (Data Source)

Use this data source to retrieve an Okta Policy.

## Example Usage

```hcl
data "okta_policy" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
