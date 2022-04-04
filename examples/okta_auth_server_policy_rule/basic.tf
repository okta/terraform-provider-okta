data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_inline_hook" "test" {
  name    = "testAcc_replace_with_uuid"
  version = "1.0.1"
  type    = "com.okta.oauth2.tokens.transform"

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test"
    method  = "POST"
  }

  auth = {
    key   = "Authorization"
    type  = "HEADER"
    value = "123"
  }
}

resource "okta_auth_server_policy_rule" "test" {
  auth_server_id       = okta_auth_server.test.id
  policy_id            = okta_auth_server_policy.test.id
  status               = "ACTIVE"
  name                 = "test"
  priority             = 1
  group_whitelist      = [data.okta_group.all.id]
  grant_type_whitelist = ["implicit"]
  inline_hook_id       = okta_inline_hook.test.id
}

resource "okta_auth_server" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}

resource "okta_auth_server_policy" "test" {
  name             = "test"
  description      = "test"
  priority         = 1
  client_whitelist = ["ALL_CLIENTS"]
  auth_server_id   = okta_auth_server.test.id
}
