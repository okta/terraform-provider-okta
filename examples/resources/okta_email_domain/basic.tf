resource "okta_email_domain" "test" {
  brand_id     = "bndtxk5pLfCA0gtSEao5"
  domain       = "testAcc-replace_with_uuid.example.com"
  display_name = "test"
  user_name    = "fff"
  validation_subdomain = "mail"
}