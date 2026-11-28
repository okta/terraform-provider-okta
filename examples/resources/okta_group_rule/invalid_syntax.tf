# Target group for the rule
resource "okta_group" "test1" {
  name        = "testAcc1_replace_with_uuid"
  description = "Test group for acceptance testing"
}

# Group rule with invalid syntax (unclosed parenthesis)
resource "okta_group_rule" "test1" {
  name              = "testAcc1_replace_with_uuid"
  status            = "ACTIVE"
  group_assignments = [okta_group.test1.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "(user.firstName == \"TestAcc\"" # Missing closing parenthesis
}
