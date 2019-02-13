resource "okta_auth_server" "test" {
  description = "test"
  name        = "testAcc_%[1]d"
}

data "okta_auth_server" "test" {
  name = "testAcc_%[1]d"
}
