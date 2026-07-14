resource "okta_admin_role_custom" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "Testing custom role with unsupported permission conditions"
  permissions = ["okta.apps.manage"]
  
  permission_conditions {
    permission = "okta.apps.manage"
    include    = ["appId"]
    exclude    = ["appType"]
  }
}