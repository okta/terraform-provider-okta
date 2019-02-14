resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "test-acc-%[1]d@testing.com"
  email      = "test-acc-%[1]d@testing.com"
}

data "okta_user" "test" {
  search {
    name  = "profile.firstName"
    value = "${okta_user.test.first_name}"
  }

  search {
    name  = "profile.lastName"
    value = "${okta_user.test.last_name}"
  }
}
