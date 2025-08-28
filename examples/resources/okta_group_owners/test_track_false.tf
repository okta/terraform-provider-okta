# Test fixture: track_all_owners = false, do not remove external owner

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

resource "okta_user" "external" {
  first_name = "TestAcc"
  last_name  = "External"
  login      = "testAcc-replace_with_uuid-owners-external@example.com"
  email      = "testAcc-replace_with_uuid-owners-external@example.com"
}

resource "okta_group" "grp" {
  name = "testAcc_replace_with_uuid_group_owners"
}

# Seed an extra owner outside of the resource via a separate resource (single-owner)
resource "okta_group_owner" "external_owner" {
  group_id          = okta_group.grp.id
  id_of_group_owner = okta_user.external.id
  type              = "USER"
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

  track_all_owners = false
}
