---
page_title: "Data Source: okta_authenticator_webauthn_custom_aaguids"
description: |-
  Lists all custom AAGUIDs for a WebAuthn authenticator.
---

# Data Source: okta_authenticator_webauthn_custom_aaguids

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

Lists all custom AAGUIDs (Authenticator Attestation Globally Unique Identifiers) configured for a WebAuthn authenticator. Only custom AAGUIDs that an admin has created are returned.

## Example Usage

```terraform
data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

data "okta_authenticator_webauthn_custom_aaguids" "all" {
  authenticator_id = data.okta_authenticator.webauthn.id
}
```

## Argument Reference

- `authenticator_id` - (Required) The ID of the WebAuthn authenticator.

## Attributes Reference

- `custom_aaguids` - List of custom AAGUIDs configured for this authenticator. Each element contains:
  - `aaguid` - The AAGUID identifier.
  - `name` - The product name associated with the AAGUID.
  - `authenticator_characteristics` - Properties of the custom AAGUID authenticator.
    - `fips_compliant` - Whether the authenticator meets FIPS compliance requirements.
    - `hardware_protected` - Whether the authenticator stores the private key on hardware.
    - `platform_attached` - Whether the AAGUID is built into the authenticator or is external.
