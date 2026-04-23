data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

resource "okta_authenticator_webauthn_custom_aaguid" "sample" {
  authenticator_id = data.okta_authenticator.webauthn.id
  aaguid           = "00000000-0000-0000-0000-000000000001"
  name             = "Test Key 1"

  authenticator_characteristics {
    fips_compliant     = false
    hardware_protected = true
    platform_attached  = false
  }
}
