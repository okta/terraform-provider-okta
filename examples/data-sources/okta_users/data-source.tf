
### Lookup Users by Search Criteria

data "okta_users" "example" {
  search {
    name       = "profile.company"
    value      = "Articulate"
    comparison = "sw"
  }
}

# Search for multiple users based on a raw search expression string
data "okta_users" "example" {
  search {
    expression = "profile.department eq \"Engineering\" and (created lt \"2014-01-01T00:00:00.000Z\" or status eq \"ACTIVE\")"
  }
}

### Lookup Users by Group Membership

resource "okta_group" "example" {
  name = "example-group"
}

data "okta_users" "example" {
  group_id = okta_group.example.id

  # optionally include each user's group membership
  include_groups = true

  # optionally include each user's administrator roles
  include_roles = true
}
