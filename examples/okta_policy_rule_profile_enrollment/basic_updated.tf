resource "okta_policy_profile_enrollment" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_inline_hook" "test" {
  name    = "testAcc_replace_with_uuid"
  status  = "ACTIVE"
  type    = "com.okta.user.pre-registration"
  version = "1.0.3"

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test2"
    method  = "POST"
  }
}

resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
}

resource "okta_policy_rule_profile_enrollment" "test" {
  policy_id           = okta_policy_profile_enrollment.test.id
  inline_hook_id      = okta_inline_hook.test.id
  target_group_id     = okta_group.test.id
  unknown_user_action = "REGISTER"
  email_verification  = true
  access              = "ALLOW"
  enroll_authenticators = [
    "password"
  ]
  profile_attributes {
    name     = "email"
    label    = "Email"
    required = true
  }
  profile_attributes {
    name     = "name"
    label    = "NameBig"
    required = true
  }
  profile_attributes {
    name     = "t-shirt"
    label    = "T-Shirt Size"
    required = false
  }
}