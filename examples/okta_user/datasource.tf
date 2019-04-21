resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "test-acc-replace_with_uuid@testing.com"
  email      = "test-acc-replace_with_uuid@testing.com"
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
