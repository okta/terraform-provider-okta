data "okta_policy" "test" {
  name = "Idp Discovery Policy"
  type = "IDP_DISCOVERY"
}

resource "okta_policy_rule_idp_discovery" "test" {
  policy_id           = data.okta_policy.test.id
  priority            = 1
  name                = "testAcc_replace_with_uuid"
  selection_type      = "DYNAMIC"
  provider_expression = "login.identifier.substringAfter('@')"
  network_connection  = "ANYWHERE"
  status              = "ACTIVE"
}
