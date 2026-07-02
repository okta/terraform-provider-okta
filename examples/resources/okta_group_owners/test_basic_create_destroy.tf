# Test fixture: basic create/destroy smoke test

resource "okta_user" "owner1" {
  first_name = "TestAcc"
  last_name  = "Owner1"
  login      = "testAcc-replace_with_uuid-del1@example.com"
  email      = "testAcc-replace_with_uuid-del1@example.com"
}

resource "okta_group" "grp" {
  name = "testAcc_replace_with_uuid_group_owners_del"
}

resource "okta_group_owners" "owners" {
  group_id = okta_group.grp.id

  owner {
    type = "USER"
    id   = okta_user.owner1.id
  }
}
