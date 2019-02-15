resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "test-acc-%[1]d@testing.com"
  email      = "test-acc-%[1]d@testing.com"
}
