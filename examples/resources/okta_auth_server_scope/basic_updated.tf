resource "okta_auth_server_scope" "test" {
  consent        = "REQUIRED"
  description    = "test_updated"
  name           = "test:something"
  display_name   = "test display name updated"
  optional       = true
  auth_server_id = okta_auth_server.test.id
}

resource "okta_auth_server" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}
