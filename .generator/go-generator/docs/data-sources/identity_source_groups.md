---
page_title: "okta_identity_source_groups Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Identity Source Groups.
---

# okta_identity_source_groups (Data Source)

Use this data source to retrieve an Okta Identity Source Groups.

## Example Usage

```hcl
data "okta_identity_source_groups" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
