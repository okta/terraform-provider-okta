resource "okta_auth_server" "test" {
  audiences   = ["whatever.rise.zone"]
  description = "test"
  name        = "testAcc_replace_with_uuid"
}

data "okta_auth_server" "test" {
  name = okta_auth_server.test.name
}
