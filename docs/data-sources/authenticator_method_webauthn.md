---
page_title: "Data Source: okta_authenticator_method_webauthn"
description: |-
  Reads WebAuthn authenticator method settings including AAGUID groups and passkey configuration.
---

# Data Source: okta_authenticator_method_webauthn

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

Reads the WebAuthn authenticator method settings, including AAGUID groups, user verification preferences, and passkey configuration.

## Example Usage

```terraform
data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

data "okta_authenticator_method_webauthn" "example" {
  authenticator_id = data.okta_authenticator.webauthn.id
}
```

## Argument Reference

- `authenticator_id` - (Required) The ID of the WebAuthn authenticator.

## Attributes Reference

- `status` - The status of the WebAuthn method (`ACTIVE` or `INACTIVE`).
- `user_verification` - User verification setting for enrollment.
- `user_verification_for_verify` - User verification setting for authentication.
- `attachment` - Method attachment setting.
- `enable_autofill_ui` - Whether the passkeys autofill UI is enabled.
- `resident_key_requirement` - Resident key requirement setting.
- `show_sign_in_with_a_passkey_button` - Whether the "Sign in with a Passkey" button is shown.
- `cert_based_attestation_validation` - Whether certificate-based attestation validation is enabled.
- `hardware_protected` - Whether the authenticator must store the private key on hardware.
- `fips_compliant` - Whether the authenticator must be FIPS compliant.
- `allow_syncable_passkeys` - Whether syncable passkeys are allowed.
- `rp_id` - The Relying Party (RP) ID configuration for WebAuthn. Contains:
  - `enabled` - Whether the RP ID is active and used for WebAuthn operations.
  - `domain` - The RP domain configuration. Contains:
    - `name` - The RP ID domain value used for WebAuthn operations.
    - `validation_status` - The validation status of the domain.
- `aaguid_groups` - The FIDO2 AAGUID groups. Each element contains:
  - `name` - The name of the AAGUID group.
  - `aaguids` - List of FIDO2 AAGUIDs in this group.
