resource "okta_email_domain" "example" {
  brand_id             = "abc123"
  domain               = "example.com"
  display_name         = "test"
  user_name            = "paul_atreides"
  validation_subdomain = "mail"
}
