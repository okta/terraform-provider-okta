data "okta_brands" "test" {
}

data "okta_email_template" "forgot_password" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  name     = "ForgotPassword"
}
