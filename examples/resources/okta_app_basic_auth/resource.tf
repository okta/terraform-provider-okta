resource "okta_app_basic_auth" "example" {
  label    = "Example"
  url      = "https://example.com/login.html"
  auth_url = "https://example.com/auth.html"
}
