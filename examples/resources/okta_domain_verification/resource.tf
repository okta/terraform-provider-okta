resource "okta_domain" "example" {
  name = "www.example.com"
}

resource "okta_domain_verification" "example" {
  domain_id = okta_domain.test.id
}
