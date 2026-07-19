---
page_title: "okta_principal_rate_limit Resource - terraform-provider-okta"
description: |-
  Manages an Okta Principal Rate Limit resource.
---

# okta_principal_rate_limit (Resource)

Manages an Okta Principal Rate Limit.

## Example Usage

```hcl
resource "okta_principal_rate_limit" "example" {
  principal_id = "<principal-id>"
  principal_type = "<principal_type>"

  # Optional fields
  # default_concurrency_percentage = 0
  # default_percentage = 0
}
```

## Schema

### Required

- `principal_id` (String) The unique identifier of the principal. This is the ID of the API token or OAuth 2.0 app.
- `principal_type` (String) The type of principal, either an API token or an OAuth 2.0 app

### Optional

- `default_concurrency_percentage` (Number) The default percentage of a given concurrency limit threshold that the owning principal can consume
- `default_percentage` (Number) The default percentage of a given rate limit threshold that the owning principal can consume

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created_by` (String) The Okta user ID of the user who created the principle rate limit entity
- `created_date` (String) The date and time the principle rate limit entity was created
- `last_update` (String) The date and time the principle rate limit entity was last updated
- `last_updated_by` (String) The Okta user ID of the user who last updated the principle rate limit entity
- `org_id` (String) The unique identifier of the Okta org
