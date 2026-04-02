---
page_title: "Data Source: okta_app_federated_claim"
description: |-
  Get a federated claim for an Okta application.
---

# Data Source: okta_app_federated_claim

Get a federated claim for an Okta application.

Use this data source to retrieve information about a federated claim that has been configured for an application. Federated claims add custom claims to tokens produced for an application using Okta Expression Language.

## Example Usage

```terraform
data "okta_app_federated_claim" "example" {
  app_id = "0oa1234567890abcdef"
  id     = "ofcu234567890abcdef"
}

output "claim_expression" {
  value = data.okta_app_federated_claim.example.expression
}
```

### Using with a Resource

```terraform
resource "okta_app_saml" "test_app" {
  label                    = "example"
  sso_url                  = "https://example.com"
  recipient                = "https://example.com"
  destination              = "https://example.com"
  audience                 = "https://example.com/audience"
  subject_name_id_template = "$${user.userName}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed          = true
  signature_algorithm      = "RSA_SHA256"
  digest_algorithm         = "SHA256"
}

resource "okta_app_federated_claim" "example" {
  app_id     = okta_app_saml.test_app.id
  name       = "role_last_name"
  expression = "user.profile.lastName"
}

data "okta_app_federated_claim" "example" {
  app_id = okta_app_saml.test_app.id
  id     = okta_app_federated_claim.example.id
}

output "claim_name" {
  value = data.okta_app_federated_claim.example.name
}

output "claim_expression" {
  value = data.okta_app_federated_claim.example.expression
}
```

## Schema

### Required

- `app_id` (String) The ID of the application that the federated claim belongs to.
- `id` (String) The unique identifier for the federated claim.

### Read-Only

- `expression` (String) The Okta Expression Language expression to be evaluated at runtime.
- `name` (String) The name of the claim to be used in the produced token.
