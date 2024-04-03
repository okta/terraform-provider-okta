resource "okta_policy_password" "example" {
  name                   = "example"
  status                 = "ACTIVE"
  description            = "Example"
  password_history_count = 4
  groups_included        = ["${data.okta_group.everyone.id}"]
}
