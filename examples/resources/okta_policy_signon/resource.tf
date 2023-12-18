resource "okta_policy_signon" "example" {
  name            = "example"
  status          = "ACTIVE"
  description     = "Example"
  groups_included = ["${data.okta_group.everyone.id}"]
}
