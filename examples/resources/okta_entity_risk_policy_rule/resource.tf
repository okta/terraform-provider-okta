data "okta_entity_risk_policy" "example" {
}

resource "okta_entity_risk_policy_rule" "example" {
  policy_id              = data.okta_entity_risk_policy.example.id
  name                   = "High Risk Response"
  risk_level             = "HIGH"
  terminate_all_sessions = true
}
