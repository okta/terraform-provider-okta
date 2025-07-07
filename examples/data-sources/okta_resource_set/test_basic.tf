resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "A test resource set for datasource testing"
  resources   = ["https://${var.hostname}/api/v1/users"]
}

data "okta_resource_set" "test" {
  id = okta_resource_set.test.id
} 
