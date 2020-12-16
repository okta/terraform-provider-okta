resource "okta_user_schema" "test" {
  index       = "testAcc_replace_with_uuid"
  title       = "terraform acceptance test"
  type        = "array"
  description = "terraform acceptance test"
  array_type  = "string"
  required    = false
  permissions = "READ_ONLY"
  master      = "PROFILE_MASTER"
}
