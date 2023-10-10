resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["client_credentials"]
  jwks_uri       = "https://example.com"
}

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
    format("%s/api/v1/apps", local.org_url)
  ]
}

resource "okta_app_oauth_role_assignment" "test" {
  client_id    = okta_app_oauth.test.client_id
  type         = "CUSTOM"
  role         = okta_admin_role_custom.test.id
  resource_set = okta_resource_set.test.id
}
