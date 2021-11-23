---
layout: 'okta'
page_title: 'Okta: okta_captcha_org_wide_settings'
sidebar_current: 'docs-okta-resource-captcha'
description: |-
    Manages Org-Wide CAPTCHA settings
---

# okta_captcha_org_wide_settings

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

This resource allows you to configure which parts of the authentication flow requires users to pass the CAPTCHA logic.
CAPTCHA org-wide settings can be disabled by unsetting `captcha_id` and `enabled_for`.

## Example Usage

```hcl
resource "okta_captcha" "example" {
  name       = "My CAPTCHA"
  type       = "HCAPTCHA"
  site_key   = "some_key"
  secret_key = "some_secret_key"
}

resource "okta_captcha_org_wide_settings" "example" {
  captcha_id  = okta_captcha.test.id
  enabled_for = ["SSR"]
}
```

The following example disables org-wide CAPTCHA.

```hcl
resource "okta_captcha" "example" {
  name       = "My CAPTCHA"
  type       = "HCAPTCHA"
  site_key   = "some_key"
  secret_key = "some_secret_key"
}

resource "okta_captcha_org_wide_settings" "example" {
}
```

## Argument Reference

The following arguments are supported:

- `captcha_id` (Optional) The ID of the CAPTCHA. 

- `enabled_for` (Optional) Array of pages that have CAPTCHA enabled. Valid values: `"SSR"`, `"SSPR"` and `"SIGN_IN"`.

## Import

Org-Wide CAPTCHA settings can be imported without any parameters.

```
$ terraform import okta_captcha_org_wide_settings.example _
```
