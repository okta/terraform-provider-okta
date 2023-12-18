resource "okta_group_rule" "example" {
  name   = "example"
  status = "ACTIVE"
  group_assignments = [
  "<group id>"]
  expression_type  = "urn:okta:expression:1.0"
  expression_value = "String.startsWith(user.firstName,\"andy\")"
}
