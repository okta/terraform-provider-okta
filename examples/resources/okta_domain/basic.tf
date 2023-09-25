resource "okta_domain" "test" {
  name                    = "testAcc.example.edu"
  certificate_source_type = "MANUAL"
}
