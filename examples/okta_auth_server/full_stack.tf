resource "okta_auth_server" "test" {
  name        = "testAcc_%[1]d"
  description = "Just, ya know, testing"
  audiences   = ["whatever.rise.zone"]
}

resource "okta_auth_server_claim" "test" {
  auth_server_id = "${okta_auth_server.test.id}"
  name           = "test"
  status         = "ACTIVE"
  claim_type     = "RESOURCE"
  value_type     = "EXPRESSION"
  value          = "cool"
}

resource "okta_auth_server_scope" "test" {
  auth_server_id = "${okta_auth_server.test.id}"
  consent        = "REQUIRED"
  description    = "This is a scope"
  name           = "test:something"
}

resource "okta_auth_server_policy" "test" {
  auth_server_id   = "${okta_auth_server.test.id}"
  status           = "ACTIVE"
  name             = "test"
  description      = "Policy"
  priority         = 2
  client_whitelist = ["ALL_CLIENTS"]
}

data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_auth_server_policy_rule" "test" {
  auth_server_id = "${okta_auth_server.test.id}"
  policy_id      = "${okta_auth_server_policy.test.id}"
  status         = "ACTIVE"
  name           = "test"
  priority       = 1

  assignments = {
    group_whitelist = ["${okta_group.all.id}"]
  }

  grant_type_whitelist = ["password"]
}
