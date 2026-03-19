resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_user_group_memberships" "test" {
  user_id = okta_user.test.id
  groups = [
    okta_group.test_1.id,
    okta_group.test_2.id,
  ]
}
