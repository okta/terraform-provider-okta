---
page_title: "okta_application_grants Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Application Grants.
---

# okta_application_grants (Data Source)

Use this data source to retrieve an Okta Application Grants.

## Example Usage

```hcl
data "okta_application_grants" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
