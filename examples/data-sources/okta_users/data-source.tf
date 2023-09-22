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
