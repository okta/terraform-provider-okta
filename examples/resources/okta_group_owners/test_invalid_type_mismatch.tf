# Test fixture: invalid type mismatch, pass a group id as type USER

resource "okta_group" "grp" {
  name = "testAcc_replace_with_uuid_invalid_type"
}

resource "okta_group" "owner_group" {
  name = "testAcc_replace_with_uuid_invalid_owner_group"
}

resource "okta_group_owners" "owners" {
  group_id = okta_group.grp.id

  owner {
    type = "USER"
    id   = okta_group.owner_group.id # invalid: group id with type USER
  }
}
