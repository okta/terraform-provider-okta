resource "okta_user" "user" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "test-acc-%[1]d@testing.com"
  email       = "test-acc-%[1]d@testing.com"
  status      = "ACTIVE"
}

resource "okta_group" "group" {
  name = "testAcc_%[1]d"
}

resource "okta_bookmark_app" "test" {
  label = "testAcc_%[1]d"
  url   = "https://test.com"

  users = {
    id       = "${okta_user.user.id}"
    username = "${okta_user.user.email}"
  }
}
