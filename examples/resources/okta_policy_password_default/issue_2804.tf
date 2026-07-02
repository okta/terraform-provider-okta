resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "Group for issue 2804 regression"
}

resource "okta_policy_password" "extra" {
  name              = "testAcc_replace_with_uuid_extra"
  description       = "Non-default policy to bump default priority"
  groups_included   = [okta_group.test.id]
  status            = "ACTIVE"
  priority          = 1
  question_recovery = "INACTIVE"
  call_recovery     = "INACTIVE"
  sms_recovery      = "INACTIVE"
}

resource "okta_policy_password_default" "test" {
  password_history_count = 5
  depends_on             = [okta_policy_password.extra]
}
