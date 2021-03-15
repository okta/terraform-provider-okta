---
layout: 'okta'
page_title: 'Okta: okta_policy_mfa_default'
sidebar_current: 'docs-okta-resource-policy-mfa-default'
description: |-
  Configures default MFA Policy.
---

# okta_policy_mfa_default

Configures default MFA Policy.

This resource allows you to configure default MFA Policy. 

## Example Usage

```hcl
resource "okta_policy_mfa_default" "default" {

}
```

## Argument Reference

The following arguments are supported:

- `duo` - (Optional) DUO [MFA policy settings](#mfa-settings).

- `fido_u2f` - (Optional) Fido U2F [MFA policy settings](#mfa-settings).

- `fido_webauthn` - (Optional) Fido Web Authn [MFA policy settings](#mfa-settings).

- `google_otp` - (Optional) Google OTP [MFA policy settings](#mfa-settings).

- `okta_call` - (Optional) Okta Call [MFA policy settings](#mfa-settings).

- `okta_otp` - (Optional) Okta OTP [MFA policy settings](#mfa-settings).

- `okta_password` - (Optional) Okta Password [MFA policy settings](#mfa-settings).

- `okta_push` - (Optional) Okta Push [MFA policy settings](#mfa-settings).

- `okta_question` - (Optional) Okta Question [MFA policy settings](#mfa-settings).

- `okta_sms` - (Optional) Okta SMS [MFA policy settings](#mfa-settings).
  
- `okta_email` - (Optional) Okta Email [MFA policy settings](#mfa-settings).

- `rsa_token` - (Optional) RSA Token [MFA policy settings](#mfa-settings).

- `symantec_vip` - (Optional) Symantec VIP [MFA policy settings](#mfa-settings).

- `yubikey_token` - (Optional) Yubikey Token [MFA policy settings](#mfa-settings).
  
- `hotp` - (Optional) HMAC-based One-Time Password [MFA policy settings](#mfa-settings).

### MFA Settings

All MFA settings above have the following structure.

- `enroll` - (Optional) Requirements for user initiated enrollment. Can be `"NOT_ALLOWED"`, `"OPTIONAL"`, or `"REQUIRED"`. By default, it is `"OPTIONAL"`.

- `consent_type` - (Optional) User consent type required before enrolling in the factor: `"NONE"` or `"TERMS_OF_SERVICE"`. By default, it is `"NONE"`.

## Attributes Reference

- `id` - ID of the default policy.

- `name` - Default policy name.

- `description` - Default policy description.

- `priority` - Default policy priority.

- `status` - Default policy status.

- `default_included_group_id` - ID of the default Okta group.

## Import

Default MFA Policy can be imported without providing Okta ID.

```
$ terraform import okta_policy_mfa_default.example .
```
