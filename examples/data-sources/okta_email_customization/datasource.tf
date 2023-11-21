data "okta_brands" "test" {
}

resource "okta_email_customization" "forgot_password_en" {
  brand_id         = tolist(data.okta_brands.test.brands)[0].id
  template_name    = "ForgotPassword"
  language         = "en"
  is_default       = true
  force_is_default = "create,destroy"
  subject          = "Forgot Password"
  body             = "Hi $$user.firstName,<br/><br/>Click this link to reset your password: $$resetPasswordLink"
}

data "okta_email_customizations" "forgot_password" {
  brand_id      = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
  depends_on = [
    okta_email_customization.forgot_password_en
  ]
}

data "okta_email_customization" "forgot_password_en" {
  brand_id         = tolist(data.okta_brands.test.brands)[0].id
  template_name    = "ForgotPassword"
  customization_id = tolist(data.okta_email_customizations.forgot_password.email_customizations)[0].id
}
