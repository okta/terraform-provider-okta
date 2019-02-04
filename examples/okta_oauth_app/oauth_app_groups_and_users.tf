resource "okta_group" "testAcc_group_%[1]d" {
  name = "testAcc_%[1]d"
}

resource "okta_user" "testAcc_user_%[1]d" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "test-acc-%[1]d@testing.com"
  email       = "test-acc-%[1]d@testing.com"
  status      = "ACTIVE"
}

resource "okta_oauth_app" "testAcc_%[1]d" {
  label                     = "testAcc_%[1]d"
  type                      = "web"
  grant_types               = ["implicit", "authorization_code"]
  redirect_uris             = ["http://d.com/"]
  post_logout_redirect_uris = ["http://d.com/post"]
  login_uri                 = "http://test.com"
  response_types            = ["code", "token", "id_token"]

  user = {
    id       = "${okta_user.testAcc_user_%[1]d.id}"
    username = "${okta_user.testAcc_user_%[1]d.email}"
  }

  groups = ["${okta_group.testAcc_group_%[1]d.id}"]
}
