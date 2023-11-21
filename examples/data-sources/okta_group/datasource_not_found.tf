resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
}

# Should fail to find the group since the type is the wrong type
data "okta_group" "test_type" {
  include_users = true
  name          = okta_group.test.name
  type          = "APP_GROUP"
}
