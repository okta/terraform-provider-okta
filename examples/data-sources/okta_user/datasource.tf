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

  depends_on = [
    okta_user_schema_property.test_array,
    okta_user_schema_property.test_number,
  ]
}

resource "okta_user" "other" {
  first_name = "Some"
  last_name  = "One"
  login      = "testAcc-someone-replace_with_uuid@example.com"
  email      = "testAcc-someone-replace_with_uuid@example.com"

  custom_profile_attributes = <<JSON
  {
    "${okta_user_schema_property.test_array.index}": ["cool","feature"]
  }
JSON

  depends_on = [
    okta_user_schema_property.test_array,
    okta_user_schema_property.test_number,
  ]
}

data "okta_user" "first_and_last" {
  search {
    name  = "profile.firstName"
    value = okta_user.test.first_name
  }

  search {
    name  = "profile.lastName"
    value = okta_user.test.last_name
  }

  depends_on = [
    okta_user.test,
    okta_user.other
  ]
}

data "okta_user" "read_by_id" {
  user_id = okta_user.test.id

  depends_on = [
    okta_user.test,
    okta_user.other
  ]
}

data "okta_user" "compound_search" {
  compound_search_operator = "or"

  search {
    name  = "profile.lastName"
    value = "Jones"
  }

  search {
    name  = "profile.lastName"
    value = okta_user.other.last_name
  }

  delay_read_seconds = 2

  depends_on = [
    okta_user.test,
    okta_user.other
  ]
}

data "okta_user" "expression_search" {
  search {
    expression = "profile.array123 eq \"feature\" and (created gt \"2021-01-01T00:00:00.000Z\")"
  }

  delay_read_seconds = 2

  depends_on = [
    okta_user.test,
    okta_user.other
  ]
}
