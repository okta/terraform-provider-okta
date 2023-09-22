# Get a single user by their id value
data "okta_user" "example" {
  user_id = "00u22mtxlrJ8YkzXQ357"
}

# Search for a single user based on supported profile properties
data "okta_user" "example" {
  search {
    name  = "profile.firstName"
    value = "John"
  }

  search {
    name  = "profile.lastName"
    value = "Doe"
  }
}

# Search for a single user based on a raw search expression string
data "okta_user" "example" {
  search {
    expression = "profile.firstName eq \"John\""
  }
}
