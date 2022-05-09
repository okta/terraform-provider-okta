---
layout: 'okta'
page_title: 'Okta: okta_email_sender'
sidebar_current: 'docs-okta-resource-email-sender'
description: |-
  Creates custom email sender.
---

# okta_email_sender

This resource allows you to create and configure a custom email sender.

## Example Usage

```hcl
resource "okta_email_sender" "example" {
  from_name    = "Paul Atreides"
  from_address = "no-reply@caladan.planet"
  subdomain    = "mail"
}
```

## Argument Reference

The following arguments are supported:

- `from_name` - (Required) Name of sender.

- `from_address` - (Required) Email address to send from.

- `subdomain` - (Required) Mail domain to send from.

## Attributes Reference

- `id` - ID of the sender.

- `status` - Status of the sender (shows whether the sender is verified).

- `dns_records` - TXT and CNAME records to be registered for the domain.
  - `fqdn` - DNS record name.
  - `record_type` - Record type can be TXT or CNAME.
  - `value` - DNS verification value

## Import

Custom email sender can be imported via the Okta ID.

```
$ terraform import okta_email_sender.example &#60;sender id&#62;
```
