resource "okta_policy_mfa_default" "test" {
  okta_password = {
    enroll = "REQUIRED"
  }
  depends_on = [okta_factor.okta_password]
}

resource "okta_factor" "okta_password" {
  provider_id = "okta_password"
}
