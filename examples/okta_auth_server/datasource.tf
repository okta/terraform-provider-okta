resource "okta_auth_server" "test" {
  description = "test"
  name        = "test%[1]d"
}

data "okta_auth_server" "test" {
  name = "test%[1]d"
}
