---
layout: 'okta'
page_title: 'Okta: okta_domain'
sidebar_current: 'docs-okta-resource-domain'
description: |-
  Manages custom domain for your organization.
---

# okta_domain

Manages custom domain for your organization.

## Example Usage

```hcl
resource "okta_domain" "example" {
  name = "www.example.com"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Custom Domain name.

- `certificate_source_type` - (Optional) Certificate source type that indicates whether the certificate is provided by the user or Okta. Accepted values: `MANUAL`, `OKTA_MANAGED`. Default value = `MANUAL`

  ~> **WARNING**: Use of `OKTA_MANAGED` requires a feature flag to be enabled.

## Attributes Reference

- `id` - Domain ID

- `validation_status` - Status of the domain.

- `dns_records` - TXT and CNAME records to be registered for the Domain.
  - `expiration` - TXT record expiration.
  - `fqdn` - DNS record name.
  - `record_type` - Record type can be TXT or CNAME.
  - `values` - DNS verification value

## Import

Okta Admin Role Targets can be imported via the Okta ID.

```
$ terraform import okta_domain.example &#60;domain_id&#62;
```
