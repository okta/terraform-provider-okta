resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group_rule" "test" {
  name              = "testAcc_replace_with_uuid"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"andy\")"
}

# This example is for syntax purposes only. If it was actually run
# data.okta_group_rule.test_by_id would fail because okta_group_rule.test
# wouldn't be in the search index yet. The data source implementation relies on
# a group rule search function in the Okta API

data "okta_group_rule" "test_by_id" {
  id = okta_group_rule.test.id
}

data "okta_group_rule" "test_by_name" {
  name = "testAcc_replace_with_uuid"
}
