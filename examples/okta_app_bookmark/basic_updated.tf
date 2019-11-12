resource "okta_user" "user" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "test-acc-replace_with_uuid@example.com"
  email       = "test-acc-replace_with_uuid@example.com"
}

resource "okta_group" "group" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_app_bookmark" "test" {
  label = "testAcc_replace_with_uuid"
  url   = "https://test.com"

  users {
    id       = "${okta_user.user.id}"
    username = "${okta_user.user.email}"
  }

  groups = ["${okta_group.group.id}"]
}
