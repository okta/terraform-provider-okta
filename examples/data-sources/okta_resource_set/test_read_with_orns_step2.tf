resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "A test resource set with multiple ORN references"
  resources_orn = [
    "orn:okta:directory:${var.orgID}:groups",
    "orn:okta:directory:${var.orgID}:apps",
    "orn:okta:directory:${var.orgID}:users"
  ]
}

data "okta_resource_set" "test" {
  id = okta_resource_set.test.id
}
