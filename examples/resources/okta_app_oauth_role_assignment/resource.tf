# Standard Role:

resource "okta_app_oauth" "test" {
  label          = "test"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["client_credentials"]
  jwks_uri       = "https://example.com"
}

resource "okta_app_oauth_role_assignment" "test" {
  client_id = okta_app_oauth.test.client_id
  type      = "HELP_DESK_ADMIN"
}

# Custom Role:

resource "okta_app_oauth" "test" {
  label          = "test"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["client_credentials"]
  jwks_uri       = "https://example.com"
}

resource "okta_admin_role_custom" "test" {
  label       = "test"
  description = "testing, testing"
  permissions = ["okta.apps.assignment.manage", "okta.users.manage", "okta.apps.manage"]
}

resource "okta_resource_set" "test" {
  label       = "test"
  description = "testing, testing"
  resources = [
    format("%s/api/v1/users", "https://example.okta.com"),
    format("%s/api/v1/apps", "https://example.okta.com")
  ]
}

resource "okta_app_oauth_role_assignment" "test" {
  client_id    = okta_app_oauth.test.client_id
  type         = "CUSTOM"
  role         = okta_admin_role_custom.test.id
  resource_set = okta_resource_set.test.id
}
