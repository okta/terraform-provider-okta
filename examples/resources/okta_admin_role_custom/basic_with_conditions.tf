resource "okta_admin_role_custom" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "Testing custom role with permission conditions"
  permissions = ["okta.users.read"]

  permission_conditions {
    permission = "okta.users.read"
    include = jsonencode({
      "okta:ResourceAttribute/User/Profile" = ["department", "costCenter"]
    })
  }
}

resource "okta_admin_role_custom" "test1" {
  label       = "testAcc_replace_with_uuid_1"
  description = "Testing custom role with permission conditions"
  permissions = ["okta.users.read"]

  permission_conditions {
    permission = "okta.users.read"
    exclude = jsonencode({
      "okta:ResourceAttribute/User/Profile" = ["title"]
    })
  }
}
