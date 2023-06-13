data "okta_brands" "test" {
}

resource "okta_email_domain" "test" {
  brand_id     = tolist(data.okta_brands.test.brands)[0].id
  domain       = "example.com"
  display_name = "test"
  user_name    = "fff"
}