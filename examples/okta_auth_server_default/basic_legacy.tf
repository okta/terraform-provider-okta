resource "okta_auth_server_default" "sun_also_rises" {
  name        = "default"
  audiences   = ["api://default"]
  description = "Default Authorization Server for your Applications"
}
