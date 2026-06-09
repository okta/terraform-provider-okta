data "okta_brands" "test" {
}

resource "okta_email_domain" "test" {
  brand_id             = tolist(data.okta_brands.test.brands)[0].id
  domain               = "testAcc-replace_with_uuid.example.com"
  display_name         = "test"
  user_name            = "fff"
  validation_subdomain = "mail"
}
