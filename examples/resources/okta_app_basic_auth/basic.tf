resource "okta_app_basic_auth" "test" {
  label                   = "testAcc_replace_with_uuid"
  url                     = "https://example.com/login.html"
  auth_url                = "https://example.com/auth.html"
  logo                    = "../examples/resources/okta_app_basic_auth/terraform_icon.png"
  reveal_password         = true
  credentials_scheme      = "EDIT_USERNAME_AND_PASSWORD"
  user_name_template      = "user.firstName"
  user_name_template_type = "CUSTOM"
  admin_note              = "admin note"
  enduser_note            = "enduser note"
}
