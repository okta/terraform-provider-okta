resource "okta_user_subscription" "test" {
  user_id           = "usr00000000000001"
  notification_type = "APP_IMPORT"
}
