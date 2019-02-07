resource "okta_user" "testAcc_%[1]d" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "test-acc-%[1]d@testing.com"
  email      = "test-acc-%[1]d@testing.com"
}

data "okta_user" "test" {
  search {
    name  = "profile.firstName"
    value = "${okta_user.testAcc_%[1]d.first_name}"
  }

  search {
    name  = "profile.lastName"
    value = "${okta_user.testAcc_%[1]d.last_name}"
  }
}

data "okta_user" "test2" {
  search {
    name  = "profile.login"
    value = "${okta_user.testAcc_%[1]d.login}"
  }
}
