data "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group_rule" "inval" {
  name              = "testAcc_replace_with_uuid"
  status            = "ACTIVE"
  group_assignments = [data.okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,String.toLowerCase(\"bob\"))"
}
