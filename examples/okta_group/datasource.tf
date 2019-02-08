resource "okta_group" "testAcc_%[1]d" {
  name        = "something new"
  description = "testing, testing"
}

data "okta_group" "testAcc_%[1]d" {
  name = "${okta_group.testAcc_%[1]d.name}"
}
