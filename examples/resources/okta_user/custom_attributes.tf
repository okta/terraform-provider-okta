resource "okta_user_schema_property" "test" {
  index  = "customAttribute123"
  title  = "terraform acceptance test"
  type   = "string"
  master = "PROFILE_MASTER"
}

resource "okta_user_schema_property" "test_array" {
  index      = "array123"
  title      = "terraform acceptance test"
  type       = "array"
  array_type = "string"
  master     = "PROFILE_MASTER"
  depends_on = [okta_user_schema_property.test_number, okta_user_schema_property.test]
}

resource "okta_user_schema_property" "test_number" {
  index      = "number123"
  title      = "terraform acceptance test"
  type       = "number"
  master     = "PROFILE_MASTER"
  depends_on = [okta_user_schema_property.test]
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"

  custom_profile_attributes = <<JSON
  {
    "customAttribute123": "testing-custom-attribute"
  }
JSON

  depends_on = [okta_user_schema_property.test, okta_user_schema_property.test_number]
}
