---
page_title: "okta_custom_domain Resource - terraform-provider-okta"
description: |-
  Manages an Okta Custom Domain resource.
---

# okta_custom_domain (Resource)

Manages an Okta Custom Domain.

## Example Usage

```hcl
resource "okta_custom_domain" "example" {

  # Optional fields
  # brand_id = "<brand-id>"
  # certificate_source_type = "<certificate_source_type>"
  # expiration = "<expiration>"
  # fqdn = "<fqdn>"
  # record_type = "<record_type>"
}
```

## Schema

### Optional

- `brand_id` (String) The ID number of the brand
- `certificate_source_type` (String) Certificate source type that indicates whether the certificate is provided by the user or Okta.
- `expiration` (String) DNS TXT record expiration
- `fqdn` (String) DNS record name
- `record_type` (String) RecordType
- `values` (List) DNS record value
- `domain` (String) Custom domain name
- `expiration` (String) Certificate expiration
- `fingerprint` (String) Certificate fingerprint
- `subject` (String) Certificate subject
- `validation_status` (String) Status of the domain

### Read-Only

- `id` (String) The unique identifier for the resource.
