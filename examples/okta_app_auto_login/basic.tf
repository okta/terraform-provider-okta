resource "okta_app_auto_login" "test" {
  label                   = "testAcc_replace_with_uuid"
  sign_on_url             = "https://example.com/login.html"
  sign_on_redirect_url    = "https://example.com"
  reveal_password         = true
  credentials_scheme      = "EDIT_USERNAME_AND_PASSWORD"
  user_name_template      = "user.firstName"
  user_name_template_type = "CUSTOM"
}
