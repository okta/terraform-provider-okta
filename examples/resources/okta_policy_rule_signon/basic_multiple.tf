data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_network_zone" "test" {
  name     = "testAcc_replace_with_uuid"
  type     = "IP"
  gateways = ["34.82.0.0/15"]
  status   = "ACTIVE"
}

resource "okta_policy_signon" "test" {
  name            = "testAcc_replace_with_uuid"
  status          = "ACTIVE"
  description     = "Terraform Acceptance Test SignOn Policy"
  groups_included = [data.okta_group.all.id]
}

resource "okta_policy_rule_signon" "test_allow" {
  policy_id          = okta_policy_signon.test.id
  name               = "testAccAllow_replace_with_uuid"
  status             = "ACTIVE"
  network_connection = "ZONE"
  network_includes   = [okta_network_zone.test.id]
}

resource "okta_policy_rule_signon" "test_deny" {
  policy_id          = okta_policy_signon.test.id
  name               = "testAccDeny_replace_with_uuid"
  status             = "ACTIVE"
  access             = "DENY"
  network_connection = "ZONE"
  network_excludes   = [okta_network_zone.test.id]
}
