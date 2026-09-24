resource "okta_role_subscription" "test" {
  role_ref          = "SUPER_ADMIN"
  notification_type = "APP_IMPORT"
}
