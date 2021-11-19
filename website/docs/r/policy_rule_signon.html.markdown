---
layout: 'okta'
page_title: 'Okta: okta_policy_rule_signon'
sidebar_current: 'docs-okta-resource-policy-rule-signon'
description: |-
  Creates a Sign On Policy Rule.
---

# okta_policy_rule_signon

Creates a Sign On Policy Rule. In case `Invalid condition type specified: riskScore.` error is thrown, set `risc_level`
to an empty string, since this feature is not enabled.

## Example Usage

```hcl
resource "okta_policy_signon" "test" {
  name = "Example Policy"
  status = "ACTIVE"
  description = "Example Policy"
}

data "okta_behavior" "new_city" {
  name = "New City"
}

resource "okta_policy_rule_signon" "example" {
  access = "CHALLENGE"
  authtype = "RADIUS"
  name = "Example Policy Rule"
  network_connection = "ANYWHERE"
  policy_id = okta_policy_signon.example.id
  status = "ACTIVE"
  risc_level = "HIGH"
  behaviors = [data.okta_behavior.new_city.id]
  factor_sequence {
    primary_criteria_factor_type = "token:hotp" // TOTP
    primary_criteria_provider = "CUSTOM"
    secondary_criteria {
      factor_type = "token:software:totp" // Okta Verify
      provider = "OKTA"
    }
    secondary_criteria { // Okta Verify Push
      factor_type = "push"
      provider = "OKTA"
    }
    secondary_criteria { // Password
      factor_type = "password"
      provider = "OKTA"
    }
    secondary_criteria { // Security Question
      factor_type = "question"
      provider = "OKTA"
    }
    secondary_criteria { // SMS
      factor_type = "sms"
      provider = "OKTA"
    }
    secondary_criteria { // Google Auth
      factor_type = "token:software:totp"
      provider = "GOOGLE"
    }
    secondary_criteria { // Email
      factor_type = "email"
      provider = "OKTA"
    }
    secondary_criteria { // Voice Call
      factor_type = "call"
      provider = "OKTA"
    }
    secondary_criteria { // FIDO2 (WebAuthn)
      factor_type = "webauthn"
      provider = "FIDO"
    }
    secondary_criteria { // RSA
      factor_type = "token"
      provider = "RSA"
    }
    secondary_criteria { // Symantec VIP
      factor_type = "token"
      provider = "SYMANTEC"
    }
  }
  factor_sequence {
    primary_criteria_factor_type = "token:software:totp" // Okta Verify
    primary_criteria_provider = "OKTA"
  }
}
```

## Argument Reference

The following arguments are supported:

- `policyid` - (Deprecated) Policy ID.
  
- `policy_id` - (Required) Policy ID.

- `name` - (Required) Policy Rule Name.

- `priority` - (Optional) Policy Rule Priority, this attribute can be set to a valid priority. To avoid endless diff situation we error if an invalid priority is provided. API defaults it to the last (lowest) if not there.

- `status` - (Optional) Policy Rule Status: `"ACTIVE"` or `"INACTIVE"`.

- `authtype` - (Optional) Authentication entrypoint: `"ANY"`, `"LDAP_INTERFACE"` or `"RADIUS"`.

- `access` - (Optional) Allow or deny access based on the rule conditions: `"ALLOW"`, `"DENY"` or `"CHALLENGE"`. The default is `"ALLOW"`.

- `mfa_required` - (Optional) Require MFA. By default is `false`.

- `mfa_prompt` - (Optional) Prompt for MFA based on the device used, a factor session lifetime, or every sign-on attempt: `"DEVICE"`, `"SESSION"` or `"ALWAYS"`.

- `mfa_remember_device` - (Optional) Remember MFA device. The default `false`.

- `mfa_lifetime` - (Optional) Elapsed time before the next MFA challenge.

- `session_idle` - (Optional) Max minutes a session can be idle.,

- `session_lifetime` - (Optional) Max minutes a session is active: Disable = 0.

- `session_persistent` - (Optional) Whether session cookies will last across browser sessions. Okta Administrators can never have persistent session cookies.

- `network_connection` - (Optional) Network selection mode: `"ANYWHERE"`, `"ZONE"`, `"ON_NETWORK"`, or `"OFF_NETWORK"`.

- `network_includes` - (Optional) The network zones to include. Conflicts with `network_excludes`.

- `network_excludes` - (Optional) The network zones to exclude. Conflicts with `network_includes`.

- `risc_level` - (Optional) Risc level: `"ANY"`, `"LOW"`, `"MEDIUM"` or `"HIGH"`. Default is `"ANY"`. It can be also 
  set to an empty string in case `RISC_SCORING` org feature flag is disabled.

- `behaviors` - (Optional) List of behavior IDs.

- `factor_sequence` - (Optional) Auth factor sequences. Should be set if `access = "CHALLENGE"`.
  - `primary_criteria_provider` - (Required) Primary provider of the auth section.
  - `primary_criteria_factor_type` - (Required) Primary factor type of the auth section.
  - `secondary_criteria` - (Optional) Additional authentication steps.
    - `provider` - (Required) Provider of the additional authentication step.
    - `factor_type` - (Required) Factor type of the additional authentication step.

- `primary_factor` - (Optional) Rule's primary factor. **WARNING** Ony works as a part of the Identity Engine. Valid values: 
`"PASSWORD_IDP_ANY_FACTOR"`, `"PASSWORD_IDP"`.

- `users_excluded` - (Optional) The list of user IDs that would be excluded when rules are processed.

## Attributes Reference

- `id` - ID of the Rule.

- `policyid` - (Deprecated) Policy ID.
  
- `policy_id` - Policy ID.

## Import

A Policy Rule can be imported via the Policy and Rule ID.

```
$ terraform import okta_policy_rule_signon.example <policy id>/<rule id>
```
