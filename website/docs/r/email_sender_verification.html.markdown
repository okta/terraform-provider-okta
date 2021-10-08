---
layout: 'okta' 
page_title: 'Okta: okta_email_sender_verification' 
sidebar_current: 'docs-okta-resource-email-sender-verification'
description: |-
  Verifies the email sender.
---

# okta_email_sender_verification

Verifies the email sender. The resource won't be created if the email sender could not be verified.

## Example Usage

```hcl
resource "okta_email_sender" "example" {
  from_name    = "Paul Atreides"
  from_address = "no-reply@caladan.planet"
  subdomain    = "mail"
}

resource "okta_email_sender_verification" "example" {
  sender_id = okta_email_sender.valid.id
}
```

## Argument Reference

The following arguments are supported:

- `sender_id` - (Required) Email sender ID.

## Import

This resource does not support importing.
