resource "okta_app_swa" "test" {
  label          = "testAcc_replace_with_uuid"
  status         = "INACTIVE"
  button_field   = "btn-login-updated"
  password_field = "txtbox-password-updated"
  username_field = "txtbox-username-updated"
  url            = "https://example.com/login-updated.html"
}
