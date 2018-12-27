resource "okta_group" "testAcc_group_%[1]d" {
  name        = "TestACC-%[1]d"
  description = "An acceptance test created group"
}

resource "okta_user" "testAcc_%[1]d" {
  first_name  = "TestAcc"
  last_name   = "Smith"
  login       = "test-acc-%[1]d@testing.com"
  email       = "test-acc-%[1]d@testing.com"
}
