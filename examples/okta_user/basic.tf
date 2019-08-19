resource "okta_user" "testAcc_replace_with_uuid" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "test-acc-replace_with_uuid@example.com"
  email      = "test-acc-replace_with_uuid@example.com"
}
