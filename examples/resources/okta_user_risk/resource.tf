resource "okta_user_risk" "example" {
  user_id    = okta_user.example.id
  risk_level = "HIGH"
}
