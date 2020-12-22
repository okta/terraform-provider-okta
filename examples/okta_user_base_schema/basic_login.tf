resource "okta_user_base_schema" "login" {
  index    = "login"
  title    = "Username"
  type     = "string"
  pattern  = "[a-z]+"
  required = true
}
