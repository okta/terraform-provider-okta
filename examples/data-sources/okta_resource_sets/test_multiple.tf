variable "hostname" {
  type = string
}

resource "okta_resource_set" "test1" {
  label       = "testAcc_replace_with_uuid_1"
  description = "First test resource set for datasource testing"
  resources   = ["https://${var.hostname}/api/v1/users"]
}

resource "okta_resource_set" "test2" {
  label       = "testAcc_replace_with_uuid_2"
  description = "Second test resource set for datasource testing"
  resources   = ["https://${var.hostname}/api/v1/groups"]

  depends_on = [okta_resource_set.test1]
}

data "okta_resource_sets" "test" {
  depends_on = [okta_resource_set.test1, okta_resource_set.test2]
}
