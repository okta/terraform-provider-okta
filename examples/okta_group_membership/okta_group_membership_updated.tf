resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Bould"
  login      = "steve_replace_with_uuid@ledzeppelin.com"
  email      = "steve_replace_with_uuid@ledzeppelin.com"

  lifecycle {
    ignore_changes = [group_memberships]
  }
}

resource "okta_group_membership" "test" {
  group_id = okta_group.test.id
  user_id  = okta_user.test.id
}
