resource "okta_group" "test3" {
  name        = "testAcc3_replace_with_uuid"
  description = "Test group for acceptance testing"
}

resource "okta_group_rule" "test3" {
  name              = "testAcc3_replace_with_uuid"
  status            = "ACTIVE"
  group_assignments = [okta_group.test3.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "user.firstName == \"TestAcc\" AND" # Trailing operator
}
