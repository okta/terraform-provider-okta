resource "okta_group" "test" {
  name        = "testAcc"
  description = "testing, testing"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Jones"
  login      = "john_replace_with_uuid@ledzeppelin.com"
  email      = "john_replace_with_uuid@ledzeppelin.com"
}

resource "okta_user_group_memberships" "test" {
  user_id = okta_user.test.id
  groups = [
    okta_group.test.id,
  ]
}

data "okta_group" "test" {
  include_users = true
  name          = okta_group.test.name

  depends_on = [ okta_user_group_memberships.test ]
}