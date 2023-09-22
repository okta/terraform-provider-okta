resource "okta_user_schema_property" "test" {
  index  = "customAttribute123"
  title  = "terraform acceptance test"
  type   = "string"
  master = "PROFILE_MASTER"
}

resource "okta_user_schema_property" "test2" {
  index  = "customAttribute1234"
  title  = "terraform acceptance test"
  type   = "string"
  master = "PROFILE_MASTER"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"

  custom_profile_attributes_to_ignore = ["customAttribute123"]
  custom_profile_attributes           = <<JSON
  {
    "customAttribute1234": "testing-custom-attribute"
  }
JSON

  depends_on = [okta_user_schema_property.test, okta_user_schema_property.test2]
}
