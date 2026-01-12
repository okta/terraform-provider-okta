resource "okta_authenticator" "webauthn" {
  name   = "WebAuthn"
  key    = "webauthn"
  status = "ACTIVE"
  provider_json = jsonencode({
    "type" : "FIDO",
    "configuration" : {}
  })
}
