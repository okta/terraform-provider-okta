resource "okta_rate_limit_admin_notification_settings" test{
  notifications_enabled = true
}

data "okta_rate_limit_admin_notification_settings" "test" {
}