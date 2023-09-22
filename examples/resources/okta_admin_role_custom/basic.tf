resource "okta_admin_role_custom" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing, testing"
  permissions = ["okta.apps.assignment.manage"]
}
