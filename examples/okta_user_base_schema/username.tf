resource "okta_user_base_schema" "login" {
  index       = "login"
  title       = "Username"
  type        = "string"
  master      = "PROFILE_MASTER"
  permissions = "READ_ONLY"
  required    = true
  pattern     = "[a-z0-9]+"
}
