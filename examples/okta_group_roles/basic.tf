resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}

resource "okta_group_roles" "test" {
  group_id    = okta_group.test.id
  admin_roles = ["SUPER_ADMIN"]
}

resource "okta_user" "test" {
  first_name        = "TestAcc"
  last_name         = "Smith"
  login             = "testAcc-replace_with_uuid@example.com"
  email             = "testAcc-replace_with_uuid@example.com"
  group_memberships = [okta_group.test.id]
}
