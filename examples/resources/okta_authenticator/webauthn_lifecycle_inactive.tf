resource "okta_authenticator" "webauthn" {
  name   = "WebAuthn"
  key    = "webauthn"
  status = "INACTIVE"
}
