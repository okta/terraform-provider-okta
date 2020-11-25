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
  policyid = okta_policy_signon.test.id
  name     = "testAcc_replace_with_uuid"
  status   = "ACTIVE"
}
