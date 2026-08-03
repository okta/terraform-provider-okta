resource "okta_admin_role_custom" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "workflow permission alias migration"
  permissions = ["okta.workflows.flows.read"]
}
