resource "okta_group" "test" {
  name = "test_423_response_capture"
}

# Main group rule that should succeed
resource "okta_group_rule" "test" {
  name              = "test_423_response_capture"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"andy\")"
}

# Additional rules that might trigger 423s but won't fail the test
resource "okta_group_rule" "concurrent1" {
  name              = "test_423_response_capture_concurrent1"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.lastName,\"smith\")"
}

resource "okta_group_rule" "concurrent2" {
  name              = "test_423_response_capture_concurrent2"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.email,\"test\")"
}

resource "okta_group_rule" "concurrent3" {
  name              = "test_423_response_capture_concurrent3"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.department,\"engineering\")"
}

resource "okta_group_rule" "concurrent4" {
  name              = "test_423_response_capture_concurrent4"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.title,\"developer\")"
}

resource "okta_group_rule" "concurrent5" {
  name              = "test_423_response_capture_concurrent5"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.employeeNumber,\"123\")"
}
