resource "okta_group" "test_other" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group_rule" "test" {
  name              = "testAcc_replace_with_uuid"
  status            = "INACTIVE"
  group_assignments = [okta_group.test_other.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"bob\")"
}
