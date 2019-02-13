resource "okta_auth_server_scope" "test" {
  consent        = "REQUIRED"
  description    = "This is a scope, part deux"
  name           = "test:something"
  auth_server_id = "${okta_auth_server.test.id}"
}

resource "okta_auth_server" "test" {
  name        = "test%[1]d"
  description = "test"
  audiences   = ["api://default"]
}
