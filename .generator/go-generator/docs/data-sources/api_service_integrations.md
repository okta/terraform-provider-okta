---
page_title: "okta_api_service_integrations Data Source - terraform-provider-okta"
description: |-
  Use this data source to retrieve an Okta Api Service Integrations.
---

# okta_api_service_integrations (Data Source)

Use this data source to retrieve an Okta Api Service Integrations.

## Example Usage

```hcl
data "okta_api_service_integrations" "example" {
  id = "<resource-id>"
}
```

## Schema

### Required

- `id` (String) ID of the resource to look up.
