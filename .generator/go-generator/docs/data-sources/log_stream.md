---
page_title: "okta_log_stream Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Log Stream.
---

# okta_log_stream (Data Source)

Use this data source to retrieve an Okta Log Stream.

## Example Usage

```hcl
data "okta_log_stream" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
