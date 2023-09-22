resource "okta_user_base_schema_property" "firstName" {
  index       = "firstName"
  permissions = "READ_ONLY"
  title       = "First name"
  type        = "string"
}
