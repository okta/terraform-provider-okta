resource "okta_security_events_provider" "example" {
  name       = "Security Events Provider"
  type       = "okta"
  is_enabled = "ACTIVE"
  settings {
    issuer   = "https://example.oktapreview.com"
    jwks_url = "https://example.oktapreview.com/oauth2/v1/keys"
  }
}
data "okta_security_events_provider" "example" {
  id = okta_security_events_provider.example.id
}