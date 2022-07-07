resource "okta_policy_mfa_default" "test" {
  okta_password = {
    enroll = "REQUIRED"
  }
  google_otp = {
    enroll = "OPTIONAL"
  }
  depends_on = [okta_factor.okta_password, okta_factor.google_otp]
}

resource "okta_factor" "okta_password" {
  provider_id = "okta_password"
}

resource "okta_factor" "google_otp" {
  provider_id = "google_otp"
}
