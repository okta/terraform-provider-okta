resource "okta_domain" "test" {
  name                    = "testAcc-replace_with_uuid.example.com"
  certificate_source_type = "MANUAL"
}
