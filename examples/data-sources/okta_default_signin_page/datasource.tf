data "okta_brands" "test" {
}


data "okta_default_signin_page" "test" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
}
