resource "okta_security_notification_emails" "example" {
  report_suspicious_activity_enabled       = true
  send_email_for_factor_enrollment_enabled = true
  send_email_for_factor_reset_enabled      = true
  send_email_for_new_device_enabled        = true
  send_email_for_password_changed_enabled  = true
}
