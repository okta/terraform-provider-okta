resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "native"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code"]
}

resource "okta_app_user_schema_property" "test" {
  app_id      = okta_app_oauth.test.id
  index       = "testAcc_replace_with_uuid"
  title       = "terraform acceptance test"
  type        = "array"
  description = "testing"
  required    = false
  permissions = "READ_ONLY"
  master      = "PROFILE_MASTER"
  array_type  = "object"
  array_enum = [
    jsonencode({ value = "test_value_1" }),
    jsonencode({ value = "test_value_2" })
  ]
  array_one_of {
    const = jsonencode({ value = "test_value_1" })
    title = "object 1"
  }
  array_one_of {
    const = jsonencode({ value = "test_value_2" })
    title = "object 2"
  }
}