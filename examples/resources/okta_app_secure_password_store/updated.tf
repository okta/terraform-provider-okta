resource "okta_app_secure_password_store" "test" {
  label              = "testAcc_replace_with_uuid"
  status             = "INACTIVE"
  username_field     = "user"
  password_field     = "pass"
  url                = "http://test.com"
  credentials_scheme = "EXTERNAL_PASSWORD_SYNC"
}
