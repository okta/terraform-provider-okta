variable "hostname" {
  type = string
}

resource "okta_group" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testing"
}

resource "okta_resource_set" "test" {
  label       = "test"
  description = "testing, testing"
  resources = [
    format("https://%s/api/v1/users", var.hostname)
  ]
}

resource "okta_admin_role_custom" "test" {
  label       = "testt"
  description = "testing, testing"
  permissions = ["okta.apps.assignment.manage", "okta.users.manage", "okta.apps.manage"]
}

resource "okta_group_role" "test" {
  group_id        = okta_group.test.id
  role_id         = okta_admin_role_custom.test.id
  resource_set_id = okta_resource_set.test.id
  role_type       = "CUSTOM"
}
