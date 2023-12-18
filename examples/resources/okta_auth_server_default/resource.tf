resource "okta_auth_server_default" "example" {
  audiences   = ["api://default"]
  description = "Default Authorization Server for your Applications"
}
