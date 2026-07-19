---
page_title: "okta_role_e_custom Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Role E Custom.
---

# okta_role_e_custom (Data Source)

Use this data source to retrieve an Okta Role E Custom.

## Example Usage

```hcl
data "okta_role_e_custom" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
