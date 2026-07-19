---
page_title: "okta_push_provider Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Push Provider.
---

# okta_push_provider (Data Source)

Use this data source to retrieve an Okta Push Provider.

## Example Usage

```hcl
data "okta_push_provider" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
