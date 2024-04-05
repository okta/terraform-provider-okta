resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
}

resource "okta_group_memberships" "test" {
  group_id = okta_group.test.id
  users = [
    okta_user.test1.id,
    okta_user.test2.id,
  ]
}
