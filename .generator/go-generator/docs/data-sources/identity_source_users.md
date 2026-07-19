---
page_title: "okta_identity_source_users Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Identity Source Users.
---

# okta_identity_source_users (Data Source)

Use this data source to retrieve an Okta Identity Source Users.

## Example Usage

```hcl
data "okta_identity_source_users" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
