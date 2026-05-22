resource "okta_auth_server" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "test"
  audiences   = ["whatever.rise.zone"]
}

resource "okta_auth_server_policy" "first" {
  auth_server_id   = okta_auth_server.test.id
  status           = "ACTIVE"
  name             = "first"
  description      = "first policy"
  priority         = 1
  client_whitelist = ["ALL_CLIENTS"]

  lifecycle {
    ignore_changes = [priority]
  }
}

resource "okta_auth_server_policy" "second" {
  auth_server_id   = okta_auth_server.test.id
  status           = "ACTIVE"
  name             = "second"
  description      = "second policy"
  priority         = 2
  client_whitelist = ["ALL_CLIENTS"]

  lifecycle {
    ignore_changes = [priority]
  }
}

resource "okta_auth_server_policy" "third" {
  auth_server_id   = okta_auth_server.test.id
  status           = "ACTIVE"
  name             = "third"
  description      = "third policy"
  priority         = 3
  client_whitelist = ["ALL_CLIENTS"]

  lifecycle {
    ignore_changes = [priority]
  }
}

resource "okta_auth_server_policy_priority" "test" {
  auth_server_id = okta_auth_server.test.id
  priorities = [
    okta_auth_server_policy.first.id,
    okta_auth_server_policy.second.id,
    okta_auth_server_policy.third.id,
  ]
}
