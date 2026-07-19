resource "okta_group_role_assignment_custom" "example" {
  group_id = "<group-id>"
  type = "<type>"

  # Optional fields
  # assignment_type = "<assignment_type>"
  # status = "ACTIVE"
}
