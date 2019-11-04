resource "okta_auth_server" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "test_updated"
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
  description      = "update"
  priority         = 1
  client_whitelist = ["${okta_app_oauth.test.id}"]
}

data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_auth_server_policy_rule" "test" {
  auth_server_id       = "${okta_auth_server.test.id}"
  policy_id            = "${okta_auth_server_policy.test.id}"
  status               = "ACTIVE"
  name                 = "test"
  priority             = 1
  group_whitelist      = ["${data.okta_group.all.id}"]
  grant_type_whitelist = ["password", "implicit"]
}

resource "okta_app_oauth" "test" {
  status         = "ACTIVE"
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["https://localhost:8443/redirect_uri/"]
  response_types = ["code", "token", "id_token"]
}
