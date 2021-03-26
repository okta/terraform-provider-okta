resource "okta_user" "test" {
  admin_roles = ["APP_ADMIN", "GROUP_MEMBERSHIP_ADMIN", "HELP_DESK_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "blah"
  login       = "testAcc_replace_with_uuid@example.com"
  email       = "testAcc_replace_with_uuid@example.com"
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

resource "okta_admin_role_targets" "test_app" {
  user_id   = okta_user.test.id
  role_type = tolist(okta_user.test.admin_roles)[0]
  apps      = [format("%s.%s", okta_app_swa.test.name, okta_app_swa.test.id)]
}

resource "okta_admin_role_targets" "test_group" {
  user_id   = okta_user.test.id
  role_type = tolist(okta_user.test.admin_roles)[1]
  groups    = [okta_group.test.id]
}
