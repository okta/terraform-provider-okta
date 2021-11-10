resource "okta_auth_server_claim_default" "test" {
  name                    = "address"
  auth_server_id          = okta_auth_server.test.id
  always_include_in_token = true
}

resource "okta_auth_server" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}
