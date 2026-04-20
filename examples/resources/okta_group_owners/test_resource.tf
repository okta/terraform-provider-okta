# Test fixture for okta_group_owners

resource "okta_user" "owner1" {
  first_name = "TestAcc"
  last_name  = "Owner1"
  login      = "testAcc-replace_with_uuid-owners1@example.com"
  email      = "testAcc-replace_with_uuid-owners1@example.com"
}

resource "okta_user" "owner2" {
  first_name = "TestAcc"
  last_name  = "Owner2"
  login      = "testAcc-replace_with_uuid-owners2@example.com"
  email      = "testAcc-replace_with_uuid-owners2@example.com"
}

resource "okta_user" "owner3" {
  first_name = "TestAcc"
  last_name  = "Owner3"
  login      = "testAcc-replace_with_uuid-owners3@example.com"
  email      = "testAcc-replace_with_uuid-owners3@example.com"
}

resource "okta_group" "grp" {
  name = "testAcc_replace_with_uuid_group_owners"
}

resource "okta_group_owners" "owners" {
  group_id = okta_group.grp.id

  owner {
    type = "USER"
    id   = okta_user.owner1.id
  }
  owner {
    type = "USER"
    id   = okta_user.owner2.id
  }
}
