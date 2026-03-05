resource "okta_admin_role_custom" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "test custom role"
  permissions = ["okta.apps.assignment.manage"]
}

data "okta_admin_role_custom" "test" {
  id = okta_admin_role_custom.test.id
}
