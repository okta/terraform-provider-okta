resource "okta_user" "testAcc_%[1]d" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "test-acc-%[1]d@testing.com"
  email      = "test-acc-%[1]d@testing.com"
}
