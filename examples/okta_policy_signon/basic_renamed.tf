data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_signon" "test" {
  name            = "testAccUpdated_replace_with_uuid"
  status          = "ACTIVE"
  description     = "Terraform Acceptance Test SignOn Policy"
  groups_included = [okta_group.test.id]
}

resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}
