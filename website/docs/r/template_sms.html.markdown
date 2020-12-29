---
layout: 'okta'
page_title: 'Okta: okta_template_sms'
sidebar_current: 'docs-okta-resource-template-sms'
description: |-
  Creates an Okta SMS Template.
---

# okta_template_sms

Creates an Okta SMS Template.

This resource allows you to create and configure an Okta SMS Template.

## Example Usage

```hcl
resource "okta_template_sms" "example" {
  type = "SMS_VERIFY_CODE"
  template = "Your $${org.name} code is: $${code}"
  translations {
    language = "en"
    template = "Your $${org.name} code is: $${code}"
  }

  translations {
    language = "es"
    template = "Tu código de $${org.name} es: $${code}."
  }
}
```

## Argument Reference

The following arguments are supported:

- `type` - (Required) SMS template type

- `template` - (Required) Default SMS message

- `translations` - (Required) Set of translations for a particular template.
  - `language` - (Required) The language to map the template to.
  - `template` - (Required) The SMS message.

## Attributes Reference

- `id` - ID of the SMS Template.

## Import

An Okta SMS Template can be imported via the template type.

```
$ terraform import okta_template_sms.example <template type>
```
