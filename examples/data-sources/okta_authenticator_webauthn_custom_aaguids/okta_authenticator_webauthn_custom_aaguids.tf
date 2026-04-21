data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

data "okta_authenticator_webauthn_custom_aaguids" "sample" {
  authenticator_id = data.okta_authenticator.webauthn.id
}
