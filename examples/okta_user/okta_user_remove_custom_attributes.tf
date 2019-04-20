resource "okta_user_schema" "testAcc_schema_uuid" {
  index  = "customAttribute123"
  title  = "terraform acceptance test"
  type   = "string"
  master = "PROFILE_MASTER"
}

resource "okta_user" "testAcc_replace_with_uuid" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "Smith"
  login       = "test-acc-replace_with_uuid@testing.com"
  email       = "test-acc-replace_with_uuid@testing.com"
}
