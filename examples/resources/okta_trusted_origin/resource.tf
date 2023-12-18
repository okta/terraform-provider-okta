resource "okta_trusted_origin" "example" {
  name   = "example"
  origin = "https://example.com"
  scopes = ["CORS"]
}
