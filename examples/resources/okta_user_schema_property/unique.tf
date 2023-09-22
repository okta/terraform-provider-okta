resource "okta_user_schema_property" "testAcc_replace_with_uuid" {
  index       = "testAcc_replace_with_uuid"
  title       = "terraform acceptance test setting unique attribute to UNIQUE_VALIDATED 007"
  type        = "string"
  description = "terraform acceptance test setting unique attribute to UNIQUE_VALIDATED 007"
  required    = false
  min_length  = 1
  max_length  = 70
  permissions = "READ_WRITE"
  master      = "OKTA"
  unique      = "UNIQUE_VALIDATED"
}
