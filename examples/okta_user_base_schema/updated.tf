resource "okta_user_base_schema" "firstName" {
  index       = "firstName"
  title       = "First name"
  type        = "string"

  permissions = "READ_WRITE"
}
