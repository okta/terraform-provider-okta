data "okta_brands" "test" {
}

data "okta_email_templates" "test" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
}
