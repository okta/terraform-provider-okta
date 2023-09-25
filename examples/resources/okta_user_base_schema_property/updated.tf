resource "okta_user_base_schema_property" "firstName" {
  index       = "firstName"
  title       = "First name"
  type        = "string"
  master      = "PROFILE_MASTER"
  permissions = "READ_WRITE"
  required    = true
}
