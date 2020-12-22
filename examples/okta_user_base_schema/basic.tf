resource "okta_user_base_schema" "firstName" {
  index       = "firstName"
  permissions = "READ_ONLY"
  title       = "First name"
  type        = "string"
}
