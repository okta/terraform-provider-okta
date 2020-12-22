resource "okta_group" "test" {
  name        = "TestACC-replace_with_uuid"
  description = "An acceptance test created group"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}
