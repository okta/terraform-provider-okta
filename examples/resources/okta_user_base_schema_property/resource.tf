resource "okta_user_base_schema_property" "example" {
  index     = "firstName"
  title     = "First name"
  type      = "string"
  required  = true
  master    = "OKTA"
  user_type = data.okta_user_type.example.id
}
