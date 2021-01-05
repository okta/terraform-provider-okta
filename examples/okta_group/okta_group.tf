resource "okta_group" "test" {
  name        = "testAcc"
  description = "testing, testing"
  users       = [okta_user.test.id]
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Jones"
  login      = "john_replace_with_uuid@ledzeppelin.com"
  email      = "john_replace_with_uuid@ledzeppelin.com"
}