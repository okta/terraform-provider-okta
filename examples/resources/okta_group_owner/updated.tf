resource "okta_group" "test1" {
  name = "testAcc_replace_with_uuid_1"
}

resource "okta_group_owner" "test1" {
  group_id                  = okta_group.test1.id
  type                      = "USER"
}
