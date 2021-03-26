resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_policy_signon" "test" {
  name            = "testAcc_replace_with_uuid"
  status          = "INACTIVE"
  description     = "Terraform Acceptance Test SignOn Policy Updated"
  groups_included = [okta_group.test.id]
}
