resource "okta_auth_server_policy_rule" "test" {
  auth_server_id = "${okta_auth_server.test.id}"
  policy_id      = "${okta_auth_server_policy.test.id}"
  status         = "ACTIVE"
  name           = "test"
  description    = "test"
  priority       = 2
}

resource "okta_auth_server" "test" {
  name        = "test%[1]d"
  description = "test"
  audiences   = ["api://default"]
}

resource "okta_auth_server_policy" "test" {
  name             = "test"
  description      = "test"
  priority         = 2
  client_whitelist = ["ALL_CLIENTS"]
  auth_server_id   = ["${okta_auth_server.test.id}"]
}
