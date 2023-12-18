resource "okta_group_role" "example" {
  group_id  = "<group id>"
  role_type = "READ_ONLY_ADMIN"
}
