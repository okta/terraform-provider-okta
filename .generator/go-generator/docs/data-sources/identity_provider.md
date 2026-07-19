---
page_title: "okta_identity_provider Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Identity Provider.
---

# okta_identity_provider (Data Source)

Use this data source to retrieve an Okta Identity Provider.

## Example Usage

```hcl
data "okta_identity_provider" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
