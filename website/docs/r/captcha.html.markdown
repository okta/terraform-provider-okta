---
layout: 'okta'
page_title: 'Okta: okta_captcha'
sidebar_current: 'docs-okta-resource-captcha'
description: |-
    Creates different types of captcha.
---

# okta_captcha

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

This resource allows you to create and configure a CAPTCHA.

## Example Usage

```hcl
resource "okta_captcha" "example" {
  name       = "My CAPTCHA"
  type       = "HCAPTCHA"
  site_key   = "some_key"
  secret_key = "some_secret_key"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the captcha.

- `type` - (Required) Type of the captcha. Valid values: `"HCAPTCHA"`, `"RECAPTCHA_V2"`.

- `site_key` - (Required) Site key issued from the CAPTCHA vendor to render a CAPTCHA on a page.

- `secret_key` - (Required) Secret key issued from the CAPTCHA vendor to perform server-side validation for a CAPTCHA token.

## Attributes Reference

- `id` - ID of the captcha.

## Import

Behavior can be imported via the Okta ID.

```
$ terraform import okta_captcha.example <captcha id>
```
