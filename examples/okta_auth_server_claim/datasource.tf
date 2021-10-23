resource "okta_auth_server" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}

data "okta_auth_server_claim" "test" {
  auth_server_id = okta_auth_server.test.id
  name           = "birthdate"
}
