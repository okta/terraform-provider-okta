resource "okta_factor" "okta_email" {
  provider_id = "okta_email"
}

resource "okta_factor" "okta_password" {
  provider_id = "okta_password"
}

resource "okta_policy_mfa_default" "test" {
  is_oie = true

  okta_password = {
    enroll = "REQUIRED"
  }

  okta_email = {
    enroll = "REQUIRED"
  }

  depends_on = [okta_factor.okta_email, okta_factor.okta_password]
}
