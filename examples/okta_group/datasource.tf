resource "okta_group" "test" {
  name        = "something new"
  description = "testing, testing"
}

data "okta_group" "test" {
  name = "${okta_group.test.name}"
}
