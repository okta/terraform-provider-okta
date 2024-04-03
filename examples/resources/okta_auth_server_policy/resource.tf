resource "okta_auth_server_policy" "example" {
  auth_server_id   = "<auth server id>"
  status           = "ACTIVE"
  name             = "example"
  description      = "example"
  priority         = 1
  client_whitelist = ["ALL_CLIENTS"]
}
