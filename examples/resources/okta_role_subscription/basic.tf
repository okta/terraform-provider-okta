resource "okta_role_subscription" "test" {
  notification_type = "APP_IMPORT"
  role_type         = "SUPER_ADMIN"
  status            = "unsubscribed"
}
