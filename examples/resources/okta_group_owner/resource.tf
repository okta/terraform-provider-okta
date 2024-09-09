resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group_owner" "test" {
  group_id                  = okta_group.test.id
  id_of_group_owner         = okta_user.test.id
  type                      = "USER"
}
