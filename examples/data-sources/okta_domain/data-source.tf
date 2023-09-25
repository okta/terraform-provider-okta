resource "okta_domain" "example" {
  name = "www.example.com"
}

data "okta_domain" "by-name" {
  domain_id_or_name = "www.example.com"

  depends_on = [
    okta_domain.example
  ]
}

data "okta_domain" "by-id" {
  domain_id_or_name = okta_domain.example.id
}
