data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

resource "okta_authenticator_method_webauthn" "sample" {
  authenticator_id  = data.okta_authenticator.webauthn.id
  user_verification = "PREFERRED"
  attachment        = "ANY"

  aaguid_group {
    name    = "TestYubiKeys"
    aaguids = ["cb69481e-8ff7-4039-93ec-0a2729a154a8"]
  }
}
