resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Jones"
  login      = "john_replace_with_uuid@ledzeppelin.com"
  email      = "john_replace_with_uuid@ledzeppelin.com"
}

data "okta_users" "compound_search" {
  compound_search_operator = "and"

  search {
    name  = "profile.firstName"
    value = "TestAcc"
  }

  search {
    name  = "profile.lastName"
    value = "Jones"
  }

  search {
    name  = "profile.email"
    comparison = "pr"
  }

  depends_on = [ okta_user.test ]
}