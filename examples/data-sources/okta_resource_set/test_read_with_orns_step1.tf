resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "A test resource set with ORN references"
  resources_orn = [
    "orn:okta:directory:${var.orgID}:users",
    "orn:okta:directory:${var.orgID}:groups"
  ]
}

data "okta_resource_set" "test" {
  id = okta_resource_set.test.id
}
