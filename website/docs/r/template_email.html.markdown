---
layout: "okta"
page_title: "Okta: okta_template_email"
sidebar_current: "docs-okta-resource-template-email"
description: |-
  Creates an Okta Email Template.
---

# okta_template_email

Creates an Okta Email Template.

This resource allows you to create and configure an Okta Email Template.

## Example Usage

```hcl
resource "okta_template_email" "example" {
  type = "email.forgotPassword"

  translations {
    language = "en"
    subject  = "Stuff"
    template = "Hi $${user.firstName},<br/><br/>Blah blah $${resetPasswordLink}"
  }

  translations {
    language = "es"
    subject  = "Cosas"
    template = "Hola $${user.firstName},<br/><br/>Puedo ir al bano $${resetPasswordLink}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) Email template type

* `translations` - (Required) Set of translations for particular template.
  * `language` - (Required) The language to map tthe template to.
  * `subject` - (Required) The email subject line.
  * `template` - (Required) The email body.

* `default_language` - (Optional) The default language, by default is set to `"en"`.

## Attributes Reference

* `id` - ID of the Email Template.

## Import

An Okta Email Template can be imported via the template type.

```
$ terraform import okta_template_email.example <template type>
```
