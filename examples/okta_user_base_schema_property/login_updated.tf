resource "okta_user_base_schema_property" "login" {
  index       = "login"
  title       = "Username"
  type        = "string"
  required    = true
  permissions = "READ_WRITE"
}
