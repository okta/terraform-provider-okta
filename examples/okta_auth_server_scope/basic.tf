resource "okta_auth_server_scope" "test" {
  consent        = "REQUIRED"
  description    = "test"
  name           = "test:something"
  auth_server_id = "${okta_auth_server.test.id}"
}

resource "okta_auth_server" "test" {
  name        = "testAcc_%[1]d"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}
