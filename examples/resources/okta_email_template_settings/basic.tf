data "okta_brands" "test" {
}

resource "okta_email_template_settings" "test" {
  brand_id      = tolist(data.okta_brands.test.brands)[0].id
  template_name = "UserActivation"
  recipients    = "NO_USERS"
}