resource "okta_group" "test" {
  name = "Everyone"
}

resource "okta_group_rule" "test" {
  name = "testAcc_%[1]d"
  status = "INACTIVE"
  group_assignments = ["${data.okta_group.id}"]
  expression_type = "urn:okta:expression:1.0"
  expression_value = "user.role==\\\"Engineer\\\""
}
