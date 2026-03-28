# Test fixture: use a group as an owner of another group

resource "okta_group" "parent" {
  name = "testAcc_replace_with_uuid_parent_group"
}

resource "okta_group" "child" {
  name = "testAcc_replace_with_uuid_child_group"
}

resource "okta_group_owners" "owners" {
  group_id = okta_group.child.id

  owner {
    type = "GROUP"
    id   = okta_group.parent.id
  }
}
