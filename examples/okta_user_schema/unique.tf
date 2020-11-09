resource "okta_user_schema" "testAcc_replace_with_uuid" {
  index       = "testAcc_replace_with_uuid"
  title       = "terraform acceptance test setting unique attribute to UNIQUE_VALIDATED"
  type        = "string"
  description = "terraform acceptance test setting unique attribute to UNIQUE_VALIDATED"
  required    = true
  min_length  = 1
  max_length  = 70
  permissions = "READ_WRITE"
  master      = "OKTA"
  unique      = "UNIQUE_VALIDATED"
}