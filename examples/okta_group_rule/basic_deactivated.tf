resource "okta_group" "test_other" {
  name = "testAcc_new_%[1]d"
}

resource "okta_group_rule" "test" {
  name              = "testAcc_%[1]d"
  status            = "INACTIVE"
  group_assignments = ["${okta_group.test_other.id}"]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.articulateId,\"auth0|\")"
}
