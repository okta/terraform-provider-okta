// Test group & user assigned to group
resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}

resource "okta_group" "test_target" {
  name        = "test_target_Acc_replace_with_uuid"
  description = "testing"
}

resource "okta_user" "test" {
  first_name        = "TestAcc"
  last_name         = "Smith"
  login             = "testAcc-replace_with_uuid@example.com"
  email             = "testAcc-replace_with_uuid@example.com"
  group_memberships = [okta_group.test.id]
}

//Usage of role & target groups
resource "okta_group_role" "test" {
  group_id  = okta_group.test.id
  role_type = "GROUP_MEMBERSHIP_ADMIN"
}

resource "okta_group_role_group_targets" "test" {
  group_id          = okta_group_role.test.group_id
  role_id           = okta_group_role.test.id
  group_target_list = [okta_group.test_target.id]
}
