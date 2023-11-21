---
layout: 'okta'
page_title: 'Okta: okta_email_domain'
sidebar_current: 'docs-okta-resource-email-domain'
description: |-
  Creates email domain.
---

# okta_email_domain

This resource allows you to create and configure an email domain.

## Example Usage

```hcl
resource "okta_email_domain" "example" {
  brand_id     = "abc123"
  domain       = "example.com"
  display_name = "test"
  user_name    = "paul_atreides"
}
```

## Argument Reference

The following arguments are supported:

- `brand_id` - (Required) Brand id of the email domain.

- `domain` - (Required) Mail domain to send from.

- `display_name` - (Required) Display name of the email domain.

- `user_name` - (Required) User name of the email domain.

## Attributes Reference

- `id` - ID of the sender.

- `validation_status` - Status of the email domain (shows whether the domain is verified).

- `dns_validation_records` - TXT and CNAME records to be registered for the domain.
  - `fqdn` - DNS record name.
  - `record_type` - Record type can be TXT or cname.
  - `value` - DNS record value
  - `expiration ` - (Deprecated) This field has been removed in the newest go sdk version and has become noop

## Import

Custom email domain can be imported via the Okta ID.

```
$ terraform import okta_email_domain.example &#60;domain id&#62;
```
