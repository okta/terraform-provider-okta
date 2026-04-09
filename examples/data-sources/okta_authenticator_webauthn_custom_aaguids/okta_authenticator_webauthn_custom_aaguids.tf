data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

resource "okta_authenticator_webauthn_custom_aaguid" "sample" {
  authenticator_id = data.okta_authenticator.webauthn.id
  aaguid           = "cb69481e-8ff7-4039-93ec-0a2729a154a8"
  name             = "YubiKey 5 Series"

  authenticator_characteristics {
    fips_compliant     = true
    hardware_protected = true
    platform_attached  = false
  }
}

data "okta_authenticator_webauthn_custom_aaguids" "sample" {
  authenticator_id = data.okta_authenticator.webauthn.id

  depends_on = [okta_authenticator_webauthn_custom_aaguid.sample]
}
