resource "okta_service_account" "example" {
  container_orn = "<container_orn>"
  name = "Example Name"
  username = "Example Username"

  # Optional fields
  # description = "Example description"
  # owner_group_ids = "<owner_group_ids>"
  # owner_user_ids = "<owner_user_ids>"
  # password = "<password>"
}
