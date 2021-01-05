resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
  users       = [okta_user.test.id]
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Jones"
  login      = "john_replace_with_uuid@ledzeppelin.com"
  email      = "john_replace_with_uuid@ledzeppelin.com"
}

# Should fail to find the group since the type is the wrong type
data "okta_group" "test_type" {
  include_users = true
  name          = okta_group.test.name
  type          = "APP_GROUP"
}
