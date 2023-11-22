data "okta_policy" "test" {
  name = "Idp Discovery Policy"
  type = "IDP_DISCOVERY"
}

resource "okta_policy_rule_idp_discovery" "test" {
  policy_id            = data.okta_policy.test.id
  priority             = 1
  name                 = "testAcc_replace_with_uuid"
  user_identifier_type = "ATTRIBUTE"
  provider_expression  = "login.identifier.substringAfter('@')"
  selection_type       = "DYNAMIC"

  // Don't have a company schema in this account, just chosing something always there
  user_identifier_attribute = "firstName"

  user_identifier_patterns {
    match_type = "EQUALS"
    value      = "Articulate"
  }
}
