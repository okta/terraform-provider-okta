data "okta_brands" "test" {
}

resource "okta_email_customization" "forgot_password_en" {
  brand_id         = tolist(data.okta_brands.test.brands)[0].id
  template_name    = "ForgotPassword"
  language         = "en"
  is_default       = true
  force_is_default = "create,destroy"
  subject          = "Forgot Password"
  body             = "Hi $$user.firstName,<br/><br/>Blah blah $$resetPasswordLink"
}

resource "okta_email_customization" "forgot_password_es" {
  brand_id         = tolist(data.okta_brands.test.brands)[0].id
  template_name    = "ForgotPassword"
  language         = "es"
  force_is_default = "create,destroy"
  subject          = "Has olvidado tu contraseña"
  body             = "Hola $$user.firstName,<br/><br/>Puedo ir al bano $$resetPasswordLink"
  depends_on = [
    okta_email_customization.forgot_password_en,
  ]
}

data "okta_email_customizations" "forgot_password" {
  brand_id      = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
  depends_on = [
    okta_email_customization.forgot_password_en,
    okta_email_customization.forgot_password_es,
  ]
}
