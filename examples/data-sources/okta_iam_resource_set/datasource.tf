variable "hostname" {
  type = string
}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing, testing"
  resources = [
    "https://${var.hostname}/api/v1/users",
  ]
  lifecycle {
    ignore_changes = [resources]
  }
}

data "okta_iam_resource_set" "test" {
  id = okta_resource_set.test.id
}
