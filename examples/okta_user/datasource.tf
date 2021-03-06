resource "okta_user_schema_property" "test_array" {
  index      = "array123"
  title      = "terraform acceptance test"
  type       = "array"
  array_type = "string"
  master     = "PROFILE_MASTER"
}

resource "okta_user_schema_property" "test_number" {
  index      = "number123"
  title      = "terraform acceptance test"
  type       = "number"
  master     = "PROFILE_MASTER"
  depends_on = [okta_user_schema_property.test_array]
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"

  custom_profile_attributes = <<JSON
  {
    "${okta_user_schema_property.test_array.index}": ["test"],
    "${okta_user_schema_property.test_number.index}": 1
  }
JSON
}

data "okta_user" "test" {
  search {
    name  = "profile.firstName"
    value = okta_user.test.first_name
  }

  search {
    name  = "profile.lastName"
    value = okta_user.test.last_name
  }
}

data "okta_user" "read_by_id" {
  user_id = okta_user.test.id
}
