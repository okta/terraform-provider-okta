resource "okta_app_user_base_schema" "test" {
  index       = "name"
  master      = "PROFILE_MASTER"
  permissions = "READ_ONLY"
  title       = "Name"
  type        = "string"
  app_id      = okta_app_oauth.test.id
  required    = false
}

resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "native"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code"]
}
