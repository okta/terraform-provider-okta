resource "okta_user" "test1" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_group" "test1" {
  name = "testAcc_replace_with_uuid_1"
}

resource "okta_group_owner" "test1" {
  group_id                  = okta_group.test1.id
  id_of_group_owner         = okta_user.test1.id
  type                      = "USER"
}
