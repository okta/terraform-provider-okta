resource "okta_user" "owner1" {
  first_name = "Alice"
  last_name  = "Owner"
  login      = "alice-owner@example.com"
  email      = "alice-owner@example.com"
}

resource "okta_user" "owner2" {
  first_name = "Bob"
  last_name  = "Owner"
  login      = "bob-owner@example.com"
  email      = "bob-owner@example.com"
}

resource "okta_group" "grp" {
  name = "demo-group"
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
