resource "okta_domain" "test" {
  name                    = "testAcc-replace_with_uuid.example.com"
  certificate_source_type = "MANUAL"
}

resource "okta_domain_verification" "test" {
  domain_id = okta_domain.test.id
}
