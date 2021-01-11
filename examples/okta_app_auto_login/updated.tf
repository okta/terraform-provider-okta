resource "okta_app_auto_login" "test" {
  label                   = "testAcc_replace_with_uuid"
  status                  = "INACTIVE"
  sign_on_url             = "https://exampleupdate.com/login.html"
  sign_on_redirect_url    = "https://exampleupdate.com"
  reveal_password         = false
  credentials_scheme      = "SHARED_USERNAME_AND_PASSWORD"
  shared_username         = "sharedusername"
  shared_password         = "sharedpassword"
  user_name_template      = "user.firstName"
  user_name_template_type = "CUSTOM"
}
