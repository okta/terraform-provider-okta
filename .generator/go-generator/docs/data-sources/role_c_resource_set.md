---
page_title: "okta_role_c_resource_set Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Role C Resource Set.
---

# okta_role_c_resource_set (Data Source)

Use this data source to retrieve an Okta Role C Resource Set.

## Example Usage

```hcl
data "okta_role_c_resource_set" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
