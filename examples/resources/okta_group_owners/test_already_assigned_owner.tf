# Test fixture: adding an already assigned owner should not fail (HTTP 400 is ignored)

resource "okta_user" "owner1" {
  first_name = "TestAcc"
  last_name  = "Owner1"
  login      = "testAcc-replace_with_uuid-aa1@example.com"
  email      = "testAcc-replace_with_uuid-aa1@example.com"
}

resource "okta_user" "owner2" {
  first_name = "TestAcc"
  last_name  = "Owner2"
  login      = "testAcc-replace_with_uuid-aa2@example.com"
  email      = "testAcc-replace_with_uuid-aa2@example.com"
}

resource "okta_group" "grp" {
  name = "testAcc_replace_with_uuid_group_owners_aa"
}

# Assign owner1 externally first (single-owner resource)
resource "okta_group_owner" "ext" {
  group_id          = okta_group.grp.id
  id_of_group_owner = okta_user.owner1.id
  type              = "USER"
}

# Now bulk-manage owners including the already-assigned owner1
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
