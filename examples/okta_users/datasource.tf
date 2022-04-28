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
}