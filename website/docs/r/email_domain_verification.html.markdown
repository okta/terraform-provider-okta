---
layout: 'okta' 
page_title: 'Okta: okta_email_domain_verification' 
sidebar_current: 'docs-okta-resource-email-domain-verification'
description: |-
  Verifies the email domain.
---

# okta_email_domain_verification

Verifies the email domain. The resource won't be created if the email domain could not be verified.

## Example Usage

```hcl
resource "okta_email_domain" "example" {
  brand_id     = "abc123"
  domain       = "example.com"
  display_name = "test"
  user_name    = "paul_atreides"
}

resource "okta_email_domain_verification" "example" {
  email_domain_id = okta_email_domain.valid.id
}
```

## Argument Reference

The following arguments are supported:

- `email_domain_id` - (Required) Email domain ID.

## Import

This resource does not support importing.
