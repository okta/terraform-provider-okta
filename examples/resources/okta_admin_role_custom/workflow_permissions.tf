resource "okta_admin_role_custom" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "workflow permission drift check"
  permissions = ["okta.workflows.invoke"]
}
