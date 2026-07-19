---
page_title: "okta_application Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Application.
---

# okta_application (Data Source)

Use this data source to retrieve an Okta Application.

## Example Usage

```hcl
data "okta_application" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
