---
layout: 'okta'
page_title: 'Okta: okta_domain'
sidebar_current: 'docs-okta-datasource-domain'
description: |-
  Get a domain from Okta.
---

# okta_domain

Use this data source to retrieve a domain from Okta.

- https://developer.okta.com/docs/reference/api/domains/#get-domain
- https://developer.okta.com/docs/reference/api/domains/#domainresponse-object

## Example Usage

```hcl
resource "okta_domain" "example" {
  name = "www.example.com"
}

data "okta_domain" "by-name" {
	domain_id_or_name = "www.example.com"

  depends_on = [
    okta_domain.example
  ]
}

data "okta_domain" "by-id" {
	domain_id_or_name = okta_domain.example.id
}
```

## Argument Reference

The following arguments are supported:

- `domain_id_or_name` - (Required) The Okta ID of the domain or the domain name itself.

## Attributes Reference

- `id` - Domain ID
- `certificate_source_type` - Certificate source type that indicates whether the certificate is provided by the user or Okta. Values: MANUAL, OKTA_MANAGED"
- `dns_records` - TXT and CNAME records to be registered for the Domain.
  - `expiration` - TXT record expiration.
  - `fqdn` - DNS record name.
  - `record_type` - Record type can be TXT or CNAME.
  - `values` - DNS verification value
- `domain` - Domain name
- `validation_status` - Status of the domain. Values: `NOT_STARTED`, `IN_PROGRESS`, `VERIFIED`, `COMPLETED`
- `public_certificate` - Certificate metadata for the Domain
