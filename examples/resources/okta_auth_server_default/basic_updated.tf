resource "okta_auth_server_default" "sun_also_rises" {
  audiences   = ["whatever.rise.zone"]
  description = "Default Authorization Server"
  status      = "ACTIVE"
}
