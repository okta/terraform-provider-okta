data "okta_org_metadata" "org" {}

resource "okta_group" "test" {
  count = 201 # Test pagination across 200-item boundary
  name  = "testAcc_${count.index}"
}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "testing pagination with large resource set"
  resources = [
    for g in okta_group.test[*] :
    "https://${data.okta_org_metadata.org.organization}/api/v1/groups/${g.id}"
  ]
}
