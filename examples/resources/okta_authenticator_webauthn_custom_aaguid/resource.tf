# First, look up the WebAuthn authenticator
data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

# Create a custom AAGUID (non-MDS example; no certificates needed)
resource "okta_authenticator_webauthn_custom_aaguid" "custom_key" {
  authenticator_id = data.okta_authenticator.webauthn.id
  aaguid           = "00000000-0000-0000-0000-000000000001"
  name             = "Custom Security Key"

  authenticator_characteristics {
    hardware_protected = true
    platform_attached  = false
  }
}

# For MDS-registered AAGUIDs (e.g., YubiKey), attestation certificates are required:
#
# resource "okta_authenticator_webauthn_custom_aaguid" "yubikey5" {
#   authenticator_id = data.okta_authenticator.webauthn.id
#   aaguid           = "cb69481e-8ff7-4039-93ec-0a2729a154a8"
#   name             = "YubiKey 5 Series"
#
#   authenticator_characteristics {
#     fips_compliant     = true
#     hardware_protected = true
#     platform_attached  = false
#   }
#
#   attestation_root_certificate {
#     x5c = "<base64-encoded-attestation-root-certificate>"
#   }
# }
