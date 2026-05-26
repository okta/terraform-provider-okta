data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_password" "test" {
  name                                  = "testAcc_replace_with_uuid"
  status                                = "ACTIVE"
  description                           = "Terraform Acceptance Test Password Policy"
  password_history_count                = 5
  groups_included                       = [data.okta_group.all.id]
  breached_password_expire_after_days   = 0
  breached_password_logout_enabled      = false
  breached_password_delegated_workflow_id = ""
}
