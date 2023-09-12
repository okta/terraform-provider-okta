resource "okta_app_shared_credentials" "test" {
  label                            = "testAcc_replace_with_uuid"
  status                           = "ACTIVE"
  button_field                     = "btn-login-updated"
  username_field                   = "txtbox-username-updated"
  password_field                   = "txtbox-password-updated"
  url                              = "https://example.com/login-updated.html"
  redirect_url                     = "https://example.com/redirect_url"
  checkbox                         = "checkbox_red-updated"
  user_name_template               = "user.firstName"
  user_name_template_type          = "CUSTOM"
  user_name_template_suffix        = "moas"
  shared_password                  = "sharedpass22"
  shared_username                  = "sharedusername22"
  accessibility_self_service       = true
  accessibility_error_redirect_url = "https://example.com/redirect_url_1"
  // deprecated in OIE
  // https://developer.okta.com/docs/reference/api/apps/#accessibility-object
  // accessibility_login_redirect_url = "https://example.com/redirect_url_2"
  auto_submit_toolbar = true
  hide_ios            = true
  logo                = "../examples/okta_app_basic_auth/terraform_icon.png"
}
