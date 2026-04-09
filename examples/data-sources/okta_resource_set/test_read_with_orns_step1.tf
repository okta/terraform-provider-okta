data "okta_org_metadata" "_" {}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "A test resource set with ORN references"
  resources_orn = [
    "orn:okta:directory:${data.okta_org_metadata._.id}:users",
    "orn:okta:directory:${data.okta_org_metadata._.id}:groups"
  ]
}

data "okta_resource_set" "test" {
  id = okta_resource_set.test.id
}
