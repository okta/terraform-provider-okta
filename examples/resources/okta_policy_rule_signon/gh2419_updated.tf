resource "okta_policy_rule_signon" "test_risk_ONLY" {
  policy_id  = "00ppe8c4f0jC1KpxX1d7"
  name       = "test_policy_risk_ONLY"
  status     = "ACTIVE"
  risk_level = "MEDIUM"
}

resource "okta_policy_rule_signon" "test_risc_ONLY" {
  policy_id  = "00ppe8c4f0jC1KpxX1d7"
  name       = "test_policy_risc_ONLY"
  status     = "ACTIVE"
  risc_level = "HIGH"
}

resource "okta_policy_rule_signon" "test_BOTH" {
  policy_id  = "00ppe8c4f0jC1KpxX1d7"
  name       = "test_policy_BOTH"
  status     = "ACTIVE"
  risk_level = "MEDIUM"
  risc_level = "HIGH"
}
