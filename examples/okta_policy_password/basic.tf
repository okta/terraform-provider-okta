data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_password" "test" {
  name                   = "testAcc_replace_with_uuid"
  status                 = "ACTIVE"
  description            = "Terraform Acceptance Test Password Policy"
  password_history_count = 4
  groups_included        = [data.okta_group.all.id]
}
