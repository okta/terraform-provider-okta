resource "okta_authenticator" "webauthn" {
  name   = "Security Key or Biometric"
  key    = "webauthn"
  status = "ACTIVE"
  provider_json = jsonencode({
    "type" : "FIDO",
    "configuration" : {}
  })
}
