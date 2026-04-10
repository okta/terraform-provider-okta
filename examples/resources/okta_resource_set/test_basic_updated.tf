data "okta_org_metadata" "org" {}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid updated"
  description = "testing, testing updated"
  resources = [
    "https://${data.okta_org_metadata.org.organization}/api/v1/users",
    "https://${data.okta_org_metadata.org.organization}/api/v1/apps"
  ]
}
