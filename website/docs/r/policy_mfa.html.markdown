---
layout: 'okta'
page_title: 'Okta: okta_policy_mfa'
sidebar_current: 'docs-okta-resource-policy-mfa'
description: |-
  Creates an MFA Policy.
---

# okta_policy_mfa

Creates an MFA Policy.

This resource allows you to create and configure an MFA Policy.

~> Requires Org Feature Flag `OKTA_MFA_POLICY`. [Contact support](mailto:dev-inquiries@okta.com) for further information.

## Example Usage

```hcl
resource "okta_policy_mfa" "classic_example" {
  name        = "MFA Policy Classic"
  status      = "ACTIVE"
  description = "Example MFA policy using Okta Classic engine with factors."
  is_oie      = false

  okta_otp = {
    enroll = "REQUIRED"
  }

  groups_included = ["${data.okta_group.everyone.id}"]
}

resource "okta_policy_mfa" "oie_example" {
  name        = "MFA Policy OIE"
  status      = "ACTIVE"
  description = "Example MFA policy that uses Okta Identity Engine (OIE) with authenticators"
  is_oie      = true

  # The following authenticator can only be used when `is_oie` is set to true
  okta_verify = {
    enroll = "REQUIRED"
  }

  groups_included = ["${data.okta_group.everyone.id}"]
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Policy Name.

- `description` - (Optional) Policy Description.

- `priority` - (Optional) Priority of the policy.

- `status` - (Optional) Policy Status: `"ACTIVE"` or `"INACTIVE"`.

- `is_oie` - (Optional) Boolean that specifies whether to use the newer Okta Identity Engine (OIE) with policy authenticators instead of the classic engine with Factors. This value determines which of the following policy factor settings can be configured. (Default = `false`)
  ~> **WARNING:** Tenant must have the Okta Identity Engine enabled in order to use this feature.

- `groups_included` - (Optional) List of Group IDs to Include.

- `duo` - (Optional) DUO [MFA policy settings](#mfa-settings) (✓ Classic, ✓ OIE).

- `external_idp` - (Optional) External IDP [MFA policy settings](#mfa-settings) (✓ OIE).

- `fido_u2f` - (Optional) Fido U2F [MFA policy settings](#mfa-settings) (✓ Classic).

- `fido_webauthn` - (Optional) Fido Web Authn [MFA policy settings](#mfa-settings) (✓ Classic).

- `google_otp` - (Optional) Google OTP [MFA policy settings](#mfa-settings) (✓ Classic, ✓ OIE).

- `hotp` - (Optional) HMAC-based One-Time Password [MFA policy settings](#mfa-settings) (✓ Classic).

- `okta_call` - (Optional) Okta Call [MFA policy settings](#mfa-settings) (✓ Classic).

- `okta_email` - (Optional) Okta Email [MFA policy settings](#mfa-settings) (✓ Classic, ✓ OIE).

- `okta_otp` - (Optional) Okta OTP (via the Okta Verify app) [MFA policy settings](#mfa-settings) (✓ Classic).

- `okta_password` - (Optional) Okta Password [MFA policy settings](#mfa-settings) (✓ Classic, ✓ OIE).

- `okta_push` - (Optional) Okta Push [MFA policy settings](#mfa-settings) (✓ Classic).

- `okta_question` - (Optional) Okta Question [MFA policy settings](#mfa-settings) (✓ Classic).

- `okta_sms` - (Optional) Okta SMS [MFA policy settings](#mfa-settings) (✓ Classic).

- `okta_verify` - (Optional) Okta Verify [MFA policy settings](#mfa-settings) (✓ OIE).

- `onprem_mfa` - (Optional) On-Prem MFA [MFA policy settings](#mfa-settings) (✓ OIE).

- `phone_number` - (Optional) Phone Number [MFA policy settings](#mfa-settings) (✓ OIE).

- `rsa_token` - (Optional) RSA Token [MFA policy settings](#mfa-settings) (✓ Classic, ✓ OIE).

- `security_question` - (Optional) Security Question [MFA policy settings](#mfa-settings) (✓ OIE).

- `symantec_vip` - (Optional) Symantec VIP [MFA policy settings](#mfa-settings) (✓ Classic).

- `webauthn` - (Optional) FIDO2 (WebAuthn) [MFA policy settings](#mfa-settings) (✓ OIE).

- `yubikey_token` - (Optional) Yubikey Token [MFA policy settings](#mfa-settings) (✓ Classic, ✓ OIE).

### MFA Settings

All MFA settings above have the following structure.

- `enroll` - (Optional) Requirements for user initiated enrollment. Can be `"NOT_ALLOWED"`, `"OPTIONAL"`, or `"REQUIRED"`. By default, it is `"OPTIONAL"`.

- `consent_type` - (Optional) User consent type required before enrolling in the factor: `"NONE"` or `"TERMS_OF_SERVICE"`. By default, it is `"NONE"`.
  ~> **NOTE:** Only applicable when using Classic mode with Factors (not OIE). When using OIE, `consent_type` is not used

## Attributes Reference

- `id` - ID of the Policy.

## Import

An MFA Policy can be imported via the Okta ID.

```
$ terraform import okta_policy_mfa.example &#60;policy id&#62;
```
