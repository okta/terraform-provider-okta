data "okta_session_violation_policy" "example" {
}

resource "okta_session_violation_policy_rule" "example" {
  policy_id                 = data.okta_session_violation_policy.example.id
  name                      = "Session Violation Rule"
  min_risk_level            = "HIGH"
  policy_evaluation_enabled = true
}
