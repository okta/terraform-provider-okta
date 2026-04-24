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
  type        = "number"
  description = "testing"
  required    = false
  permissions = "READ_ONLY"
  master      = "PROFILE_MASTER"
  enum        = ["0.011", "0.022", "0.033"]
  one_of {
    title = "number point oh one one"
    const = "0.011"
  }
  one_of {
    title = "number point oh two two"
    const = "0.022"
  }
  one_of {
    title = "number point oh three three"
    const = "0.033"
  }
}