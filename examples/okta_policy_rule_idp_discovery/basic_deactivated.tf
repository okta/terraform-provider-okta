data "okta_policy" "test" {
  name = "Idp Discovery Policy"
  type = "IDP_DISCOVERY"
}

resource "okta_policy_rule_idp_discovery" "test" {
  status               = "INACTIVE"
  policyid             = data.okta_policy.test.id
  priority             = 1
  name                 = "testAcc_replace_with_uuid"
  user_identifier_type = "IDENTIFIER"

  user_identifier_patterns {
    match_type = "SUFFIX"
    value      = "gmail.com"
  }

  user_identifier_patterns {
    match_type = "SUFFIX"
    value      = "articulate.com"
  }
}
