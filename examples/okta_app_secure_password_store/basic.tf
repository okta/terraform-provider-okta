resource "okta_app_secure_password_store" "test" {
  label              = "testAcc_replace_with_uuid"
  username_field     = "user"
  password_field     = "pass"
  url                = "http://test.com"
  credentials_scheme = "ADMIN_SETS_CREDENTIALS"
}
