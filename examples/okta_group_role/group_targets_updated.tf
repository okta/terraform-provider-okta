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
  group_memberships = [okta_group.test.id]
}

// Test Target Groups
resource "okta_group" "test_target1" {
  name        = "testTarget1Acc_replace_with_uuid"
  description = "testing"
}

//Usage of role
resource "okta_group_role" "test" {
  group_id          = okta_group.test.id
  role_type         = "HELP_DESK_ADMIN"
  target_group_list = [okta_group.test_target1.id]
}
