resource "okta_security_events_provider" "example" {
  name       = "Security Events Provider with well-known URL"
  type       = "okta"
  is_enabled = "ACTIVE"
  settings {
    issuer   = "https://example.oktapreview.com"
    jwks_url = "https://example.oktapreview.com/oauth2/v1/keys"
  }
}