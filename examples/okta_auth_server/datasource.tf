resource "okta_auth_server" "test" {
  audiences   = ["whatever.rise.zone"]
  description = "test"
  name        = "testAcc_%[1]d"
}

data "okta_auth_server" "test" {
  name = "testAcc_%[1]d"
}
