resource "okta_policy_mfa_default" "test" {
  is_oie = true

  okta_password = {
    enroll = "REQUIRED"
  }

  okta_verify_fastpass = {
    enroll = "REQUIRED"
  }

  okta_verify_push = {
    enroll = "OPTIONAL"
  }

  okta_verify_totp = {
    enroll = "OPTIONAL"
  }
}
