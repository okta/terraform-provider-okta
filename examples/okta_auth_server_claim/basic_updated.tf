resource "okta_auth_server_claim" "test" {
  name           = "test_updated"
  status         = "INACTIVE"
  claim_type     = "RESOURCE"
  value_type     = "EXPRESSION"
  value          = "cool_updated"
  auth_server_id = "${okta_auth_server.test.id}"
}

resource "okta_auth_server" "test" {
  name        = "testAcc_%[1]d"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}
