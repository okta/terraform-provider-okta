data "okta_brands" "test" {
}

data "okta_email_customizations" "forgot_password" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
}

data "okta_email_customization" "forgot_password_en" {
  customization_id = tolist(data.okta_email_customizations.forgot_password.email_customizations)[0].id
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  template_name = "ForgotPassword"
}