data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_signon" "test" {
  name            = "testAcc_replace_with_uuid"
  status          = "ACTIVE"
  description     = "Terraform Acceptance Test SignOn Policy"
  groups_included = [data.okta_group.all.id]
}

resource "okta_policy_rule_signon" "test" {
  policyid           = okta_policy_signon.test.id
  name               = "testAcc_replace_with_uuid"
  access             = "DENY"
  status             = "ACTIVE"
  network_connection = "ZONE"
  network_excludes   = [okta_network_zone.test.id]
}

resource "okta_network_zone" "test" {
  name     = "testAcc_replace_with_uuid"
  type     = "IP"
  gateways = ["1.2.3.4/24", "2.3.4.5-2.3.4.15"]
  proxies  = ["2.2.3.4/24", "3.3.4.5-3.3.4.15"]
}
