// Test group & user assigned to group
resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}

resource "okta_user" "test" {
  first_name        = "TestAcc"
  last_name         = "Smith"
  login             = "testAcc-replace_with_uuid@example.com"
  email             = "testAcc-replace_with_uuid@example.com"
}

//Usage of role
resource "okta_group_role" "test" {
  group_id  = okta_group.test.id
  role_type = "READ_ONLY_ADMIN"
}

resource "okta_group_role" "test_app" {
  group_id  = okta_group.test.id
  role_type = "APP_ADMIN"
}
