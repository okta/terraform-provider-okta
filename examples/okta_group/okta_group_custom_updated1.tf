resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
  custom_profile_attributes = jsonencode({
    "testSchema1_replace_with_uuid" = "testing1234",
    "testSchema2_replace_with_uuid" = true,
    "testSchema3_replace_with_uuid" = 54321
  })
}

// Test Schema
resource "okta_group_schema_property" "test1" {
  index       = "testSchema1_replace_with_uuid"
  title       = "TestSchema1_replace_with_uuid"
  type        = "string"
  description = "Test string schema"
  master      = "OKTA"
}

resource "okta_group_schema_property" "test2" {
  index       = "testSchema2_replace_with_uuid"
  title       = "TestSchema2_replace_with_uuid"
  type        = "boolean"
  description = "Test bool schema"
  master      = "OKTA"
}

resource "okta_group_schema_property" "test3" {
  index       = "testSchema3_replace_with_uuid"
  title       = "TestSchema3_replace_with_uuid"
  type        = "number"
  description = "Test bool schema"
  master      = "OKTA"
}
