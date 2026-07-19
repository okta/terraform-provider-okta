---
page_title: "okta_ui_schema Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Ui Schema.
---

# okta_ui_schema (Data Source)

Use this data source to retrieve an Okta Ui Schema.

## Example Usage

```hcl
data "okta_ui_schema" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
