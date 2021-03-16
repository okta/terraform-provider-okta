resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Jones"
  login      = "john_replace_with_uuid@ledzeppelin.com"
  email      = "john_replace_with_uuid@ledzeppelin.com"

  lifecycle {
    ignore_changes = [group_memberships]
  }
}

## Test Case
resource "okta_group" "test_2" {
  name        = "testAcc_2_replace_with_uuid"
  description = "testing, testing"
}

resource "okta_group_membership" "test_2" {
  group_id = okta_group.test_2.id
  user_id  = okta_user.test.id
}
