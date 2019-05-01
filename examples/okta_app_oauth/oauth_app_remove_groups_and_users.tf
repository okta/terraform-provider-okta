resource "okta_group" "group" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_user" "user" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "test-acc-replace_with_uuid@testing.com"
  email       = "test-acc-replace_with_uuid@testing.com"
  status      = "ACTIVE"
}

resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code", "token", "id_token"]
  omit_secret    = true
}
