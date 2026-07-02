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
