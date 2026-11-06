# The catch-all rule is a system rule: it cannot be deleted on its own and is
# removed by Okta when its policy is deleted. Forget it from state rather than
# destroying it, then let the policy be destroyed.
removed {
  from = okta_policy_rule_profile_enrollment.test

  lifecycle {
    destroy = false
  }
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
