resource "okta_auth_server_scope" "test" {
  consent        = "REQUIRED"
  description    = "test"
  name           = "test:something"
  display_name   = "test display name"
  optional       = false
  auth_server_id = okta_auth_server.test.id
}

resource "okta_auth_server" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}
