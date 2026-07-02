data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

resource "okta_authenticator_method_webauthn" "sample" {
  authenticator_id  = data.okta_authenticator.webauthn.id
  user_verification = "REQUIRED"
  attachment        = "ANY"

  aaguid_group {
    name    = "TestYubiKeys"
    aaguids = ["cb69481e-8ff7-4039-93ec-0a2729a154a8"]
  }

  aaguid_group {
    name    = "TestTitanKeys"
    aaguids = ["42b4fb4a-2866-43b2-bc9c-049e6c44b3a5"]
  }
}
