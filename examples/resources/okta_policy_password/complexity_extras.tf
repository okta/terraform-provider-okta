data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_password" "test" {
  name                                          = "testAcc_replace_with_uuid"
  status                                        = "ACTIVE"
  description                                   = "Terraform Acceptance Test Password Policy Complexity"
  password_max_consecutive_repeating_characters = 3
  password_oel_statement                        = "password.value.contains('okta')"
  groups_included                               = [data.okta_group.all.id]
}
