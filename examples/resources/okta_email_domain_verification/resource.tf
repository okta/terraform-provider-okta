resource "okta_email_domain" "example" {
  brand_id     = "abc123"
  domain       = "example.com"
  display_name = "test"
  user_name    = "paul_atreides"
}

resource "okta_email_domain_verification" "example" {
  email_domain_id = okta_email_domain.valid.id
}
