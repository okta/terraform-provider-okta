---
page_title: "okta_group_owners Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Group Owners.
---

# okta_group_owners (Data Source)

Use this data source to retrieve an Okta Group Owners.

## Example Usage

```hcl
data "okta_group_owners" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
