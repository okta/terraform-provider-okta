resource "okta_admin_role_custom" "example" {
  label       = "AppAssignmentManager"
  description = "This role allows app assignment management"
  permissions = ["okta.apps.assignment.manage"]
}
