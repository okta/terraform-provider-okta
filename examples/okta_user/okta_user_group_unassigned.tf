resource "okta_group" "test" {
  name        = "TestACC-%[1]d"
  description = "An acceptance test created group"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "test-acc-%[1]d@testing.com"
  email      = "test-acc-%[1]d@testing.com"
}
