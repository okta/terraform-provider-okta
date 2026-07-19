---
page_title: "okta_user_grants Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta User Grants.
---

# okta_user_grants (Data Source)

Use this data source to retrieve an Okta User Grants.

## Example Usage

```hcl
data "okta_user_grants" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
