resource "okta_brand" "example" {
  name = "test"
}

resource "okta_domain" "test" {
  name                    = "testAcc-replace_with_uuid.example.com"
  brand_id                = resource.okta_brand.example.id
  certificate_source_type = "MANUAL"
}
