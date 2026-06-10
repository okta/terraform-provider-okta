resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group" "test_other" {
  name = "other_testAcc_replace_with_uuid"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_group_rule" "test" {
  name              = "testAcc_replace_with_uuid"
  status            = "INACTIVE"
  group_assignments = [okta_group.test_other.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"bob\")"
  users_excluded    = [okta_user.test.id]
}
