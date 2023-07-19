resource "okta_group" "test_other" {
  name = "other_testAcc_replace_with_uuid"
}

resource "okta_group_rule" "test" {
  name              = "testAcc_replace_with_uuid"
  status            = "ACTIVE"
  group_assignments = [okta_group.test_other.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,String.toLowerCase(\"bOb\"))"
}
