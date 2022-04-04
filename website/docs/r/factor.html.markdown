---
layout: 'okta'
page_title: 'Okta: okta_factor'
sidebar_current: 'docs-okta-resource-factor'
description: |-
  Allows you to manage the activation of Okta MFA methods.
---

# okta_factor

Allows you to manage the activation of Okta MFA methods.

This resource allows you to manage Okta MFA methods.

## Example Usage

```hcl
resource "okta_factor" "example" {
  provider_id = "google_otp"
}
```

## Argument Reference

The following arguments are supported:

- `provider_id` - (Required) The MFA provider name.
  Allowed values are `"duo"`, `"fido_u2f"`, `"fido_webauthn"`, `"google_otp"`, `"okta_call"`, `"okta_otp"`, `"okta_password"`, `"okta_push"`, `"okta_question"`, `"okta_sms"`, `"okta_email"`, `"rsa_token"`, `"symantec_vip"`, `"yubikey_token"`, or `"hotp"`.

- `active` - (Optional) Whether to activate the provider, by default, it is set to `true`.

## Attributes Reference

- `provider_id` - MFA provider name.
