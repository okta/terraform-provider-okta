resource "okta_domain" "test" {
  name   = "www.example.com"
  verify = false
}

data "okta_domain" "by-id" {
  domain_id_or_name = okta_domain.test.id
}

data "okta_domain" "by-name" {
  domain_id_or_name = "www.example.com"

  depends_on = [
    okta_domain.test
  ]
}
