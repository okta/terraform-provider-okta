// Test group & user assigned to group
resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}

resource "okta_user" "test" {
  first_name        = "TestAcc"
  last_name         = "Smith"
  login             = "testAcc-replace_with_uuid@example.com"
  email             = "testAcc-replace_with_uuid@example.com"
}

// Test Target App
resource "okta_app_swa" "test" {
  label          = "testAcc_replace_with_uuid"
  button_field   = "btn-login"
  password_field = "txtbox-password"
  username_field = "txtbox-username"
  url            = "https://example.com/login.html"
}

//Usage of role
resource "okta_group_role" "test" {
  group_id  = okta_group.test.id
  role_type = "HELP_DESK_ADMIN"
}

resource "okta_group_role" "test_app" {
  group_id  = okta_group.test.id
  role_type = "APP_ADMIN"
}
