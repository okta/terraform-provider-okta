resource "okta_group" "test" {
  name        = "TestACC-replace_with_uuid"
  description = "An acceptance test created group"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "test-acc-replace_with_uuid@testing.com"
  email      = "test-acc-replace_with_uuid@testing.com"

  group_memberships = ["${okta_group.test.id}"]
}
