variable "hostname" {
  type = string
}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "test resource set for data source"
  resources = [
    "https://${var.hostname}/api/v1/users"
  ]
}

data "okta_resource_set" "test" {
  id = okta_resource_set.test.id
}
