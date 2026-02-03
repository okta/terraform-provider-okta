resource "okta_authenticator" "webauthn" {
  name   = "Security Key or Biometric"
  key    = "webauthn"
  status = "INACTIVE"
  depends_on = [ okta_policy_mfa_default.default_policy ]
}

resource "okta_policy_mfa_default" "default_policy" {
  webauthn = {
    enroll = "NOT_ALLOWED"
  }
  fido_webauthn = {
    enroll = "NOT_ALLOWED"
  }
  okta_password = {
    enroll = "REQUIRED"
  }
}
