resource "okta_auth_server" "test" {
  name        = "test%[1]d"
  description = "Just, ya know, testing"
  audiences   = ["api://default"]
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

resource "okta_auth_server_policy_rule" "test" {
  auth_server_id = "${okta_auth_server.test.id}"
  policy_id      = "${okta_auth_server_policy.test.id}"
  status         = "ACTIVE"
  name           = "test"
  description    = "Policy rule"
  priority       = 2
}
