resource "okta_user" "test" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "Smith"
  login       = "test-acc-replace_with_uuid@example.com"
  email       = "test-acc-replace_with_uuid@example.com"
  password    = "Abcd1234"
  status      = "STAGED"
}
