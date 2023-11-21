resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAccUpdated-replace_with_uuid@example.com"
  email      = "testAccUpdated-replace_with_uuid@example.com"
}
