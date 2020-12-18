resource "okta_user" "test" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "Smith"
  login       = "testAcc-replace_with_uuid@example.com"
  email       = "testAcc-replace_with_uuid@example.com"
}
