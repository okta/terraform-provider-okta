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

# Non-managed owner to demonstrate track_all_owners behavior
resource "okta_user" "owner3_external" {
  first_name = "Eve"
  last_name  = "NotManaged"
  login      = "eve-external@example.com"
  email      = "eve-external@example.com"
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

  # Optional
  # track_all_owners = false
}
