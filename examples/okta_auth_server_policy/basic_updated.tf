resource "okta_auth_server_policy" "test" {
  status           = "ACTIVE"
  name             = "test_updated"
  description      = "test updated"
  priority         = 2
  client_whitelist = ["ALL_CLIENTS"]
  auth_server_id   = "${okta_auth_server.test.id}"
}

resource "okta_auth_server_policy" "test_priority" {
  status           = "ACTIVE"
  name             = "other_test"
  description      = "other test updated"
  priority         = 1
  client_whitelist = ["ALL_CLIENTS"]
  auth_server_id   = "${okta_auth_server.test.id}"
}

resource "okta_auth_server" "test" {
  name        = "testAcc_%[1]d"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}
