---
page_title: "okta_hook_key Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Hook Key.
---

# okta_hook_key (Data Source)

Use this data source to retrieve an Okta Hook Key.

## Example Usage

```hcl
data "okta_hook_key" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
