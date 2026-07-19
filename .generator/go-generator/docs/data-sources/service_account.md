---
page_title: "okta_service_account Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Service Account.
---

# okta_service_account (Data Source)

Use this data source to retrieve an Okta Service Account.

## Example Usage

```hcl
data "okta_service_account" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
