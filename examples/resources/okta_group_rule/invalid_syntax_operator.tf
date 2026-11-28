resource "okta_group" "test2" {
  name = "testAcc2_replace_with_uuid"
}

resource "okta_group_rule" "test" {
  name              = "testAcc2_replace_with_uuid"
  status            = "ACTIVE"
  group_assignments = [okta_group.test2.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "user.firstName == AND user.lastName" # Invalid operator sequence
}
