resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing, testing"
  custom_profile_attributes = jsonencode({
    "${okta_group_schema_property.test1.index}" = null,
    "${okta_group_schema_property.test2.index}" = true,
    "${okta_group_schema_property.test3.index}" = null
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
