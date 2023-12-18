locals {
  org_url = "https://mycompany.okta.com"
}

resource "okta_admin_role_custom" "test" {
  label       = "SomeUsersAndApps"
  description = "Manage apps assignments and users"
  permissions = ["okta.apps.assignment.manage", "okta.users.manage", "okta.apps.manage"]
}

resource "okta_resource_set" "test" {
  label       = "UsersWithApp"
  description = "All the users and SWA app"
  resources = [
    format("%s/api/v1/users", local.org_url),
    format("%s/api/v1/apps/%s", local.org_url, okta_app_swa.test.id)
  ]
}

// this user and group will manage the set of resources based on the permissions specified in the custom role
resource "okta_admin_role_custom_assignments" "test" {
  resource_set_id = okta_resource_set.test.id
  custom_role_id  = okta_admin_role_custom.test.id
  members = [
    format("%s/api/v1/users/%s", local.org_url, okta_user.test.id),
    format("%s/api/v1/groups/%s", local.org_url, okta_group.test.id)
  ]
}

// this user will have `CUSTOM` role assigned, but it won't appear in the `admin_roles` for that user,
// since direct assignment of custom roles is not allowed
resource "okta_user" "test" {
  first_name = "Paul"
  last_name  = "Atreides"
  login      = "no-reply@caladan.planet"
  email      = "no-reply@caladan.planet"
}

resource "okta_app_swa" "test" {
  label          = "My SWA App"
  button_field   = "btn-login"
  password_field = "txtbox-password"
  username_field = "txtbox-username"
  url            = "https://example.com/login.html"
}

resource "okta_group" "test" {
  name        = "General"
  description = "General Group"
}
