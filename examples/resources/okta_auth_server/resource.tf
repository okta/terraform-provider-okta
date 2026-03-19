resource "okta_auth_server" "example" {
  audiences   = ["api://example"]
  description = "My Example Auth Server"
  name        = "example"
  issuer_mode = "CUSTOM_URL"
  status      = "ACTIVE"
}
