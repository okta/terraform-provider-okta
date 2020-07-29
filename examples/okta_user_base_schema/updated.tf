resource "okta_user_base_schema" "firstName" {
  index       = "firstName"
  title       = "First name"
  type        = "string"
  master      = "PROFILE_MASTER"
  pattern     = "[a-z]+"
  permissions = "READ_WRITE"
  required    = true
}
