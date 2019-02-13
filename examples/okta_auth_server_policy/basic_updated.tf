resource "okta_auth_server_policy" "test" {
  status           = "INACTIVE"
  name             = "test_updated"
  description      = "Policy updated"
  priority         = 2
  client_whitelist = ["ALL_CLIENTS"]
  auth_server_id   = "${okta_auth_server.test.id}"
}

resource "okta_auth_server" "test" {
  name        = "test%[1]d"
  description = "test"
  audiences   = ["api://default"]
}
