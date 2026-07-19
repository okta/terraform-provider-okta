---
page_title: "okta_api_service_integrations Resource - terraform-provider-okta"
description: |-
  Manages an Okta Api Service Integrations resource.
---

# okta_api_service_integrations (Resource)

Manages an Okta Api Service Integrations.

## Example Usage

```hcl
resource "okta_api_service_integrations" "example" {

  # Optional fields
  # granted_scopes = "<granted_scopes>"
  # properties = "<properties>"
  # type = "<type>"
}
```

## Schema

### Optional

- `granted_scopes` (List) The list of Okta management scopes granted to the API Service Integration instance. See [Okta management OAuth 2.0 scopes](/oauth2/#okta-admin-management).
- `properties` (String) App instance properties
- `type` (String) The type of the API service integration. This string is an underscore-concatenated, lowercased API service integration name. For example, `my_api_log_integration`.

### Read-Only

- `id` (String) The unique identifier for the resource.
- `config_guide_url` (String) The URL to the API service integration configuration guide
- `created_at` (String) Timestamp when the API Service Integration instance was created
- `created_by` (String) The user ID of the API Service Integration instance creator
- `name` (String) The name of the API service integration that corresponds with the `type` property. This is the full name of the API service integration listed in the Okta Integration Network (OIN) catalog.
