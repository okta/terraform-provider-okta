data "okta_brands" "test" {
}

data "okta_email_customizations" "forgot_password" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
}

resource "okta_email_customization" "forgot_password_en_alt" {
  brand_id      = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
  language      = "cs" # Setting the language to Czech for acceptance testing
  is_default    = false # A true is_default will cause an error during create if there is already a default customization. In general, use for update.
  subject       = "Forgot Password"
  body          = "Hi $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"
}