data "okta_org_metadata" "org" {}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing, testing"
  resources_orn = [
    "orn:okta:directory:${data.okta_org_metadata.org.id}:groups",
    "orn:okta:directory:${data.okta_org_metadata.org.id}:users"
  ]
}
