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
  array_type  = "string"
  array_enum  = ["true", "false"]
  array_one_of {
    const = "true"
    title = "boolean True"
  }
  array_one_of {
    const = "false"
    title = "boolean False"
  }
}