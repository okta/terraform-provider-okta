data "okta_default_policy" "test" {
  type = "MFA_ENROLL"
}

resource "okta_policy_rule_mfa" "test" {
  policy_id = data.okta_default_policy.test.id
  name      = "testAcc_replace_with_uuid"
  status    = "INACTIVE"
  enroll    = "LOGIN"
}
