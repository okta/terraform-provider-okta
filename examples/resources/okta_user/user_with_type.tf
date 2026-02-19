resource "okta_user_type" "test" {
  name         = "testAcc_replace_with_uuid"
  display_name = "Contractor"
  description  = "A contractor user type"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"

  type {
    id = okta_user_type.test.id
  }
}