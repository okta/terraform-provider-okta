resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group_owner" "test1" {
  group_id                  = okta_group.test.id
  type                      = "USER"
}
