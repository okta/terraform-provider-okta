resource "okta_policy_password_default" "test" {
  sms_recovery           = "ACTIVE"
  password_history_count = 5
}
