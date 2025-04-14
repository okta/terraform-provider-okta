resource "okta_group" "test_group1" {
  name        = "testAcc_replace_with_uuid"
  description = "Test group for acceptance testing"
}

resource "okta_group_rule" "test" {
  name              = "testAcc_replace_with_uuid"
  status            = "ACTIVE"
  group_assignments = [okta_group.test_group1.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "(user.firstName == \"John\" AND user.lastName == \"Doe\") OR user.email == \"john@example.com\""
}
