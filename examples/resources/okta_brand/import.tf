# This example 
data "okta_brands" "test" {
}

resource "okta_brand" "example" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id

  lifecycle {
    ignore_changes = [
      agree_to_custom_privacy_policy,
      custom_privacy_policy_url,
      remove_powered_by_okta
    ]
  }
}
