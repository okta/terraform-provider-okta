resource "okta_group" "group" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_user" "user" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "test-acc-replace_with_uuid@example.com"
  email       = "test-acc-replace_with_uuid@example.com"
  status      = "ACTIVE"
}


resource "okta_app_oauth" "test" {
  label                     = "testAcc_replace_with_uuid"
  type                      = "web"
  grant_types               = ["implicit", "authorization_code"]
  redirect_uris             = ["http://d.com/"]
  post_logout_redirect_uris = ["http://d.com/post"]
  login_uri                 = "http://test.com"
  response_types            = ["code", "token", "id_token"]
  consent_method            = "TRUSTED"
  implicit_assignment       = true

  // This block will be ignored when implicit_assignment is true,
  // Added in so that the acceptance test would fail due to okta 400 error should that ingore code not be working
  users {
    id       = "${okta_user.user.id}"
    username = "${okta_user.user.email}"
  }

  // This block will be ignored when implicit_assignment is true
  groups = ["${okta_group.group.id}"]
}
