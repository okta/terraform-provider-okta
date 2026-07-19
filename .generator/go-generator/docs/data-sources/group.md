---
page_title: "okta_group Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Group.
---

# okta_group (Data Source)

Use this data source to retrieve an Okta Group.

## Example Usage

```hcl
data "okta_group" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
