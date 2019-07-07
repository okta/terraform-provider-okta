resource okta_user_schema test_array {
  index      = "array123"
  title      = "terraform acceptance test"
  type       = "array"
  array_type = "string"
  master     = "PROFILE_MASTER"
}

resource okta_user_schema test_number {
  index  = "number123"
  title  = "terraform acceptance test"
  type   = "number"
  master = "PROFILE_MASTER"
}

resource okta_user test {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "test-acc-replace_with_uuid@testing.com"
  email      = "test-acc-replace_with_uuid@testing.com"

  custom_profile_attributes = <<JSON
  {
    "array123": ["test"],
    "number123": 1
  }
JSON

  depends_on = ["okta_user_schema.test_array", "okta_user_schema.test_number"]
}

data okta_user test {
  search {
    name  = "profile.firstName"
    value = "${okta_user.test.first_name}"
  }

  search {
    name  = "profile.lastName"
    value = "${okta_user.test.last_name}"
  }
}
