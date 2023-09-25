resource "okta_brand" "example" {
  name = "test"
}

resource "okta_domain" "test" {
  name                    = "testAccTest.example.edu"
  brand_id                = resource.okta_brand.example.id
  certificate_source_type = "OKTA_MANAGED"
}
