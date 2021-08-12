---
layout: 'okta'
page_title: 'Okta: okta_factor_totp'
sidebar_current: 'docs-okta-resource-factor-totp'
description: |-
  Allows you to manage the time-based one-time password (TOTP) factors.
---

# okta_factor_totp

Allows you to manage the time-based one-time password (TOTP) factors. A time-based one-time password (TOTP) is a
temporary passcode that is generated for user authentication. Examples of TOTP include hardware authenticators and
mobile app authenticators.

Once saved, the settings cannot be changed (except for the `name` field). Any other change would force resource
recreation.

## Example Usage

```hcl
resource "okta_factor_totp" "example" {
  name = "example"
  otp_length = 10
  hmac_algorithm = "HMacSHA256"
  time_step = 30
  clock_drift_interval = 10
  shared_secret_encoding = "hexadecimal"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The TOTP name.

- `otp_length` - (Optional) Length of the password. Default is `6`.

- `hmac_algorithm` - (Optional) - HMAC Algorithm. Valid values: `"HMacSHA1"`, `"HMacSHA256"`, `"HMacSHA512"`. Default
  is `"HMacSHA512"`.

- `time_step` - (Optional) - Time step in seconds. Valid values: `15`, `30`, `60`. Default is `15`.

- `clock_drift_interval` - (Optional) - Clock drift interval. This setting allows you to build in tolerance for any
  drift between the token's current time and the server's current time. Valid values: `3`, `5`, `10`. Default is `3`.

- `shared_secret_encoding` - (Optional) - Shared secret encoding. Valid values: `"base32"`, `"base64"`, `"hexadecimal"`.
  Default is `"base32"`.

## Attributes Reference

- `id` - ID of the TOTP factor.
