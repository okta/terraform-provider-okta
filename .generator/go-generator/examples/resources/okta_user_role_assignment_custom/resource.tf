resource "okta_user_role_assignment_custom" "example" {
  user_id = "<user-id>"
  type = "<type>"

  # Optional fields
  # assignment_type = "<assignment_type>"
  # status = "ACTIVE"
}
