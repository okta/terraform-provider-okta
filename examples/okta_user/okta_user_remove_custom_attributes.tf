resource "okta_user_schema" "testAcc_schema_%[1]d" {
  index     = "customAttribute123"
  title     = "terraform acceptance test"
  type      = "string"
}

resource "okta_user" "testAcc_%[1]d" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "Smith"
  login       = "test-acc-%[1]d@testing.com"
  email       = "test-acc-%[1]d@testing.com"
}
