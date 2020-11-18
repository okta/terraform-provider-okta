resource "okta_auth_server_claim" "test" {
  name              = "test"
  status            = "ACTIVE"
  claim_type        = "RESOURCE"
  value_type        = "GROUPS"
  group_filter_type = "EQUALS"
  value             = "Everyone"
  auth_server_id    = okta_auth_server.test.id
}

resource "okta_auth_server_claim" "test_sw" {
  name              = "test_sw"
  status            = "ACTIVE"
  claim_type        = "RESOURCE"
  value_type        = "GROUPS"
  group_filter_type = "STARTS_WITH"
  value             = "Every"
  auth_server_id    = okta_auth_server.test.id
}

resource "okta_auth_server" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}
