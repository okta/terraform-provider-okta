---
page_title: "okta_email_domain Resource - terraform-provider-okta"
description: |-
  Manages an Okta Email Domain resource.
---

# okta_email_domain (Resource)

Manages an Okta Email Domain.

## Example Usage

```hcl
resource "okta_email_domain" "example" {
  display_name = "Example Display Name"
  user_name = "Example User Name"
  brand_id = "<brand-id>"

  # Optional fields
  # fqdn = "<fqdn>"
  # record_type = "<record_type>"
  # verification_value = "<verification_value>"
  # domain = "<domain>"
  # validation_status = "ACTIVE"
}
```

## Schema

### Required

- `display_name` (String) DisplayName
- `user_name` (String) UserName
- `brand_id` (String) BrandId

### Optional

- `fqdn` (String) Fqdn
- `record_type` (String) RecordType
- `verification_value` (String) VerificationValue
- `domain` (String) Domain
- `validation_status` (String) ValidationStatus
- `validation_subdomain` (String) The subdomain for the email sender's custom mail domain

### Read-Only

- `id` (String) The unique identifier for the resource.
- `embedded` (String) Embedded
