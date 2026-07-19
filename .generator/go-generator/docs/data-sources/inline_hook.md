---
page_title: "okta_inline_hook Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Inline Hook.
---

# okta_inline_hook (Data Source)

Use this data source to retrieve an Okta Inline Hook.

## Example Usage

```hcl
data "okta_inline_hook" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
