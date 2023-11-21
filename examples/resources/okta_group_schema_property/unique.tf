resource "okta_group_schema_property" "test" {
  index       = "testAcc_replace_with_uuid"
  title       = "terraform acceptance test setting unique attribute to UNIQUE_VALIDATED 006"
  type        = "string"
  description = "terraform acceptance test setting unique attribute to UNIQUE_VALIDATED 006"
  required    = true
  min_length  = 1
  max_length  = 70
  permissions = "READ_WRITE"
  master      = "OKTA"
  unique      = "UNIQUE_VALIDATED"
}
