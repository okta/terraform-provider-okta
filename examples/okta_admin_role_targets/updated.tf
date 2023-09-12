resource "okta_user" "test" {
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "testAcc_replace_with_uuid@example.com"
  email       = "testAcc_replace_with_uuid@example.com"
}

resource "okta_user_admin_roles" "test" {
  user_id     = okta_user.test.id
  admin_roles = [
    "APP_ADMIN",
    "GROUP_MEMBERSHIP_ADMIN"
  ]
}

resource "okta_app_swa" "test" {
  label          = "testAcc_replace_with_uuid"
  button_field   = "btn-login"
  password_field = "txtbox-password"
  username_field = "txtbox-username"
  url            = "https://example.com/login.html"
}

resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}

resource "okta_group" "test_2" {
  name        = "testAcc_2_replace_with_uuid"
  description = "testing, testing"
}

resource "okta_admin_role_targets" "test_app" {
  user_id   = okta_user.test.id
  role_type = "APP_ADMIN"
  apps      = ["oidc_client", "facebook"]
}

resource "okta_admin_role_targets" "test_group" {
  user_id   = okta_user.test.id
  role_type = "GROUP_MEMBERSHIP_ADMIN"
  groups    = [okta_group.test.id, okta_group.test_2.id]
}
