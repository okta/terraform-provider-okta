resource "okta_policy_password_default" "test" {
  sms_recovery           = "INACTIVE"
  password_history_count = 0
}
