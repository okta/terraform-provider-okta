variable "hostname" {
  type = string
}

locals {
  org_url = "https://${var.hostname}"
}

resource "okta_admin_role_custom" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing, testing"
  permissions = ["okta.apps.assignment.manage", "okta.users.manage", "okta.apps.manage"]
}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing, testing"
  resources = [
    format("%s/api/v1/users", local.org_url),
    format("%s/api/v1/apps/%s", local.org_url, okta_app_swa.test.id)
  ]
}

resource "okta_admin_role_custom_assignments" "test" {
  resource_set_id = okta_resource_set.test.id
  custom_role_id  = okta_admin_role_custom.test.id
  members = [
    format("%s/api/v1/users/%s", local.org_url, okta_user.test.id),
    format("%s/api/v1/groups/%s", local.org_url, okta_group.test.id)
  ]
}

// this user will have `CUSTOM` role assigned, but it won't appear in the `admin_roles` for that user,
// since direct assignment for custom roles is not allowed
resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "blah"
  login      = "testAcc_replace_with_uuid@example.com"
  email      = "testAcc_replace_with_uuid@example.com"
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
