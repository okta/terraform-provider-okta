---
page_title: "Resource: okta_authenticator_method_webauthn"
description: |-
  Manages WebAuthn authenticator method settings including AAGUID groups and passkey configuration.
---

# Resource: okta_authenticator_method_webauthn

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

Manages the WebAuthn authenticator method settings, including AAGUID group allowlists, user verification preferences, and passkey configuration.

-> **Note:** This resource manages the settings of an existing WebAuthn authenticator method. The method itself cannot be created or deleted — it exists as part of the authenticator. On destroy, settings are reset to defaults.

## Example Usage

```terraform
data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

resource "okta_authenticator_method_webauthn" "example" {
  authenticator_id        = data.okta_authenticator.webauthn.id
  user_verification       = "PREFERRED"
  attachment              = "ANY"
  allow_syncable_passkeys = true

  aaguid_group {
    name    = "YubiKeys"
    aaguids = ["cb69481e-8ff7-4039-93ec-0a2729a154a8"]
  }
}
```

## Argument Reference

- `authenticator_id` - (Required, ForceNew) The ID of the WebAuthn authenticator.
- `user_verification` - (Optional) User verification setting for enrollment. Values: `DISCOURAGED`, `PREFERRED`, `REQUIRED`.
- `user_verification_for_verify` - (Optional) User verification setting for authentication (verification).
- `attachment` - (Optional) Method attachment setting.
- `allow_syncable_passkeys` - (Optional) Whether syncable passkeys are allowed.
- `enable_autofill_ui` - (Optional) Enables the passkeys autofill UI for WebAuthn discoverable credentials.
- `resident_key_requirement` - (Optional) Resident key requirement. Values: `REQUIRED`, `DISCOURAGED`, `PREFERRED`.
- `show_sign_in_with_a_passkey_button` - (Optional) Whether to show the "Sign in with a Passkey" button.
- `cert_based_attestation_validation` - (Optional) Whether certificate-based attestation validation is enabled.
- `hardware_protected` - (Optional) Whether the authenticator must store the private key on hardware.
- `fips_compliant` - (Optional) Whether the authenticator must be FIPS compliant.
- `aaguid_group` - (Optional) List of AAGUID group configurations.
  - `name` - (Required) A name to identify the group of FIDO2 AAGUIDs.
  - `aaguids` - (Required) A list of FIDO2 AAGUIDs in this group.

## Attributes Reference

- `id` - The authenticator ID.
- `status` - The status of the WebAuthn method (`ACTIVE` or `INACTIVE`).

## Import

The WebAuthn method can be imported using the authenticator ID:

```shell
terraform import okta_authenticator_method_webauthn.example aut1234567890
```
