resource "okta_group_role_assignment_user_admin" "example" {
  group_id = "<group-id>"
  type = "<type>"

  # Optional fields
  # catalog = "<catalog>"
  # profile = "<profile>"
  # type = "<type>"
  # assignment_type = "<assignment_type>"
  # status = "ACTIVE"
}
