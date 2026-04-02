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
  array_type  = "integer"
  array_enum  = [4, 5, 6]
  array_one_of {
    const = "4"
    title = "integer four"
  }
  array_one_of {
    const = "5"
    title = "integer five"
  }
  array_one_of {
    const = "6"
    title = "integer six"
  }
}