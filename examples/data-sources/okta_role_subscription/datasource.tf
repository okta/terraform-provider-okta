data "okta_role_subscription" "test" {
  role_ref = "SUPER_ADMIN"
  id       = "APP_IMPORT"
}
