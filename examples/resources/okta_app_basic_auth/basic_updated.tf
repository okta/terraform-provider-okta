resource "okta_app_basic_auth" "test" {
  label                          = "testAcc_replace_with_uuid"
  url                            = "https://example.com/login.html"
  auth_url                       = "https://example.org/auth.html"
  logo                           = "../examples/resources/okta_app_basic_auth/terraform_icon.png"
  reveal_password                = false
  credentials_scheme             = "SHARED_USERNAME_AND_PASSWORD"
  shared_username                = "sharedusername"
  shared_password                = "sharedpassword"
  user_name_template             = "user.firstName"
  user_name_template_type        = "CUSTOM"
  user_name_template_push_status = "PUSH"
  admin_note                     = "admin note updated"
  enduser_note                   = "enduser note updated"
}
