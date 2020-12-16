resource "okta_user_schema" "test" {
  index       = "testAcc_replace_with_uuid"
  title       = "terraform acceptance test updated"
  type        = "array"
  description = "terraform acceptance test updated"
  array_type  = "string"
  required    = true
  permissions = "READ_WRITE"
  master      = "OKTA"
}
