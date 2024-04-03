resource "okta_role_subscription" "test" {
  role_type         = "SUPER_ADMIN"
  notification_type = "APP_IMPORT"
  status            = "unsubscribed"
}
