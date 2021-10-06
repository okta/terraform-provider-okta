---
layout: 'okta' 
page_title: 'Okta: okta_domain_verification' 
sidebar_current: 'docs-okta-resource-domain-verification'
description: |-
  Verifies the Domain.
---

# okta_domain_verification

Verifies the Domain. This is replacement for the `verify` field from the `okta_domain` resource. The resource won't be 
created if the domain could not be verified. The provided will make several requests to verify the domain until 
the API returns `VERIFIED` verification status. 

## Example Usage

```hcl
resource "okta_domain" "example" {
  name = "www.example.com"
}

resource "okta_domain_verification" "example" {
  domain_id = okta_domain.test.id
}
```

## Argument Reference

The following arguments are supported:

- `domain_id` - (Required) Domain ID.

## Import

This resource does not support importing.
