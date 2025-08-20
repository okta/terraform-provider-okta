data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  email      = "testAcc_replace_with_uuid@gmail.com"
  login      = "testAcc_replace_with_uuid@gmail.com"
}

resource "okta_policy_signon" "test" {
  name            = "testAcc_replace_with_uuid"
  status          = "ACTIVE"
  description     = "Terraform Acceptance Test SignOn Policy"
  groups_included = [data.okta_group.all.id]
}

resource "okta_policy_signon" "test_two" {
  name            = "test_two"
  status          = "ACTIVE"
  description     = "Terraform Acceptance Test SignOn Policy"
  groups_included = [data.okta_group.all.id]
}

resource "okta_policy_rule_signon" "test" {
  policy_id          = okta_policy_signon.test.id
  name               = "testAcc_replace_with_uuid"
  status             = "INACTIVE"
  access             = "DENY"
  session_idle       = 240
  session_lifetime   = 240
  session_persistent = false
  users_excluded     = [okta_user.test.id]
}

resource "okta_policy_rule_signon" "test_risk_ONLY" {
  policy_id       = okta_policy_signon.test_two.id
  name            = "test_policy_risk_ONLY"
  status          = "ACTIVE"
  risk_level      = "MEDIUM"
}

resource "okta_policy_rule_signon" "test_risc_ONLY" {
  policy_id       = okta_policy_signon.test_two.id
  name            = "test_policy_risc_ONLY"
  status          = "ACTIVE"
  risc_level      = "HIGH"
} 

resource "okta_policy_rule_signon" "test_BOTH" {
  policy_id       = okta_policy_signon.test_two.id
  name            = "test_policy_BOTH"
  status          = "ACTIVE"
  risk_level      = "MEDIUM"
  risc_level      = "HIGH"
}
