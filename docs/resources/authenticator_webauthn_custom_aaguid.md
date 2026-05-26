---
page_title: "Resource: okta_authenticator_webauthn_custom_aaguid"
description: |-
  Manages a custom AAGUID for a WebAuthn authenticator.
---

# Resource: okta_authenticator_webauthn_custom_aaguid

~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

Manages a custom AAGUID (Authenticator Attestation Globally Unique Identifier) for a WebAuthn authenticator. Custom AAGUIDs allow administrators to register specific hardware security key models with the organization.

~> **NOTE:** AAGUIDs that are registered in the [FIDO Metadata Service (MDS)](https://fidoalliance.org/metadata/) with attestation certificates **require** the `attestation_root_certificate` block(s) to be provided. If omitted for an MDS-registered AAGUID, the API will return a validation error. AAGUIDs not present in the MDS do not require certificates.

## Example Usage

### Basic (custom/non-MDS AAGUID)

```terraform
data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

resource "okta_authenticator_webauthn_custom_aaguid" "custom_key" {
  authenticator_id = data.okta_authenticator.webauthn.id
  aaguid           = "00000000-0000-0000-0000-000000000001"
  name             = "Custom Security Key"

  authenticator_characteristics {
    hardware_protected = true
    platform_attached  = false
  }
}
```

### With attestation certificates (MDS-registered AAGUID)

```terraform
data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

resource "okta_authenticator_webauthn_custom_aaguid" "yubikey5" {
  authenticator_id = data.okta_authenticator.webauthn.id
  aaguid           = "cb69481e-8ff7-4039-93ec-0a2729a154a8"
  name             = "YubiKey 5 Series"

  authenticator_characteristics {
    fips_compliant     = true
    hardware_protected = true
    platform_attached  = false
  }

  attestation_root_certificate {
    x5c = "<base64-encoded-attestation-root-certificate>"
  }
}
```

## Argument Reference

- `authenticator_id` - (Required, ForceNew) The ID of the WebAuthn authenticator.
- `aaguid` - (Required, ForceNew) The Authenticator Attestation Globally Unique Identifier (AAGUID). A 128-bit identifier indicating the authenticator model.
- `name` - (Required) The product name associated with this AAGUID.
- `authenticator_characteristics` - (Optional) Properties of the custom AAGUID authenticator.
  - `fips_compliant` - (Optional) Indicates whether the authenticator meets FIPS compliance requirements.
  - `hardware_protected` - (Optional) Indicates whether the authenticator stores the private key on a hardware component.
  - `platform_attached` - (Optional) Indicates whether the custom AAGUID is built into the authenticator or is external.
- `attestation_root_certificate` - (Optional) List of attestation root certificates. **Required** for AAGUIDs that exist in the FIDO Metadata Service (MDS) with certificates.
  - `x5c` - (Required) X.509 certificate chain (base64-encoded).

## Attributes Reference

- `id` - The AAGUID identifier.
- `name` - The product name.
- `attestation_root_certificate` - Contains computed fields after creation:
  - `x5t_s256` - SHA-256 hash (thumbprint) of the X.509 certificate.
  - `issuer` - Issuer of the certificate.
  - `expiry` - Expiry date of the certificate.

## Import

A custom AAGUID can be imported using the format `authenticator_id/aaguid`:

```shell
terraform import okta_authenticator_webauthn_custom_aaguid.example aut1234567890/cb69481e-8ff7-4039-93ec-0a2729a154a8
```
