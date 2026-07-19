---
page_title: "okta_user_type Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta User Type.
---

# okta_user_type (Data Source)

Use this data source to retrieve an Okta User Type.

## Example Usage

```hcl
data "okta_user_type" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
