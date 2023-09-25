resource "okta_user_schema_property" "test" {
  index       = "testAcc_replace_with_uuid"
  title       = "terraform acceptance test updated 005"
  type        = "array"
  description = "terraform acceptance test updated 005"
  array_type  = "string"
  required    = false
  permissions = "READ_WRITE"
  master      = "OKTA"
}
