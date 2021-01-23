resource "okta_auth_server_claim" "test" {
  name           = "test"
  status         = "ACTIVE"
  claim_type     = "RESOURCE"
  value_type     = "EXPRESSION"
  value          = "cool"
  auth_server_id = okta_auth_server.test.id
}

resource "okta_auth_server" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}
