resource "okta_user_role_assignment_api_access_management_admin" "example" {
  user_id = "<user-id>"
  type = "<type>"

  # Optional fields
  # catalog = "<catalog>"
  # profile = "<profile>"
  # type = "<type>"
  # assignment_type = "<assignment_type>"
  # status = "ACTIVE"
}
