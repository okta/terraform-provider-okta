---
page_title: "okta_behavior Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Behavior.
---

# okta_behavior (Data Source)

Use this data source to retrieve an Okta Behavior.

## Example Usage

```hcl
data "okta_behavior" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
