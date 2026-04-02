---
page_title: "Resource: okta_app_federated_claim"
description: |-
  Manages a federated claim for an Okta application.
  Federated claims allow you to pass user information from Okta to your app integrations.
---

# Resource: okta_app_federated_claim

Manages a federated claim for an Okta application.

Federated claims allow you to pass user information from Okta to your app integrations.

## Example Usage

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
```

### Multiple Claims Example

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

resource "okta_app_federated_claim" "last_name" {
  app_id     = okta_app_saml.test_app.id
  name       = "lastName"
  expression = "user.profile.lastName"
}

resource "okta_app_federated_claim" "first_name" {
  app_id     = okta_app_saml.test_app.id
  name       = "firstName"
  expression = "user.profile.firstName"
}

resource "okta_app_federated_claim" "department" {
  app_id     = okta_app_saml.test_app.id
  name       = "department"
  expression = "user.profile.department"
}
```

## Schema

### Required

- `app_id` (String) The ID of the application to add the federated claim to.
- `expression` (String) The Okta Expression Language expression to be evaluated at runtime. See [Okta Expression Language](https://developer.okta.com/docs/reference/okta-expression-language/) for more information.
- `name` (String) The name of the claim to be used in the produced token.

### Read-Only

- `id` (String) The unique identifier for the federated claim. This is a combination of `app_id` and the claim ID separated by a forward slash (`/`).

## Import

An app federated claim can be imported using the format `app_id/id`:

```shell
terraform import okta_app_federated_claim.example <app_id>/<id>
```

Example:

```shell
terraform import okta_app_federated_claim.example 0oa1234567890abcdef/clm1234567890abcdef
```
