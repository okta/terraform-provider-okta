---
page_title: "okta_application_federated_claim Resource - terraform-provider-okta"
description: |-
  Manages an Okta Application Federated Claim resource.
---

# okta_application_federated_claim (Resource)

Manages an Okta Application Federated Claim.

## Example Usage

```hcl
resource "okta_application_federated_claim" "example" {
  app_id = "<app-id>"

  # Optional fields
  # expression = "<expression>"
  # name = "Example Name"
}
```

## Schema

### Required

- `app_id` (String) ID of the parent application

### Optional

- `expression` (String) The Okta Expression Language expression to be evaluated at runtime
- `name` (String) The name of the claim to be used in the produced token

### Read-Only

- `id` (String) The unique identifier for the resource.

## Import

Import using `{app_id}/{id}`:

```shell
terraform import okta_application_federated_claim.example <app_id> <id>
```
