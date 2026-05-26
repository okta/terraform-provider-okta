data "okta_authenticator" "webauthn" {
  key = "webauthn"
}

data "okta_authenticator_method_webauthn" "sample" {
  authenticator_id = data.okta_authenticator.webauthn.id
}
