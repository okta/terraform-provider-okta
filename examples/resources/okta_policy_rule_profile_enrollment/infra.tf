resource "okta_policy_profile_enrollment" "test" {
  name   = "testAcc_replace_with_uuid"
  status = "ACTIVE"
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
