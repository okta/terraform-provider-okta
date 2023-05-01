resource "okta_group" "test" {
  name        = "testAcc"
  description = "testing, testing"
}

data "okta_group" "test" {
  include_users = true
  name          = okta_group.test.name
}
