resource "okta_group" "group" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_user" "user" {
  first_name = "TestAcc"
  last_name  = "blah"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
  status     = "ACTIVE"
}

resource "okta_app_oauth" "test" {
  label                     = "testAcc_replace_with_uuid"
  type                      = "web"
  grant_types               = ["implicit", "authorization_code"]
  redirect_uris             = ["http://d.com/"]
  post_logout_redirect_uris = ["http://d.com/post"]
  login_uri                 = "http://test.com"
  response_types            = ["code", "token", "id_token"]
}
