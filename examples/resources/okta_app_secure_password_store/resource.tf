resource "okta_app_secure_password_store" "example" {
  label              = "example"
  username_field     = "user"
  password_field     = "pass"
  url                = "https://test.com"
  credentials_scheme = "ADMIN_SETS_CREDENTIALS"
}
