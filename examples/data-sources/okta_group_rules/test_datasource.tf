resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_group_rule" "test1" {
  name              = "testAccOne_replace_with_uuid"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"andy\")"
}

resource "okta_group_rule" "test2" {
  name              = "testAccTwo_replace_with_uuid"
  status            = "ACTIVE"
  group_assignments = [okta_group.test.id]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.firstName,\"bob\")"
  depends_on        = [okta_group_rule.test1]
}

# List all group rules
data "okta_group_rules" "all" {
  depends_on = [okta_group_rule.test1, okta_group_rule.test2]
}

# List group rules with search filter
data "okta_group_rules" "filtered" {
  search     = "testAcc"
  depends_on = [okta_group_rule.test1, okta_group_rule.test2]
}

# List group rules with limit
data "okta_group_rules" "limited" {
  limit      = 50
  depends_on = [okta_group_rule.test1, okta_group_rule.test2]
}

# List group rules with expand parameter
data "okta_group_rules" "with_expand" {
  expand     = "groupIdToGroupNameMap"
  depends_on = [okta_group_rule.test1, okta_group_rule.test2]
}
