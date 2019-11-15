resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "test-acc-replace_with_uuid@example.com"
  email      = "test-acc-replace_with_uuid@example.com"
  password   = "Abcd1234"
}
