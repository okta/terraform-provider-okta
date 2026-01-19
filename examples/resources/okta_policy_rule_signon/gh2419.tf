data "okta_group" "all" {
  name = "Everyone"
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

resource "okta_policy_rule_signon" "test_risk_ONLY" {
  policy_id  = "00ppe8c4f0jC1KpxX1d7"
  name       = "test_policy_risk_ONLY"
  status     = "ACTIVE"
  risk_level = "ANY"
}

resource "okta_policy_rule_signon" "test_risc_ONLY" {
  policy_id  = "00ppe8c4f0jC1KpxX1d7"
  name       = "test_policy_risc_ONLY"
  status     = "ACTIVE"
  risc_level = "MEDIUM"
}

resource "okta_policy_rule_signon" "test_BOTH" {
  policy_id  = "00ppe8c4f0jC1KpxX1d7"
  name       = "test_policy_BOTH"
  status     = "ACTIVE"
  risk_level = "LOW"
  risc_level = "HIGH"
}

