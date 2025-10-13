resource "okta_domain" "test-downcase" {
  name = "testdowncase.example.com"
}

data "okta_domain" "by-id-downcase" {
  domain_id_or_name = okta_domain.test-downcase.id
}