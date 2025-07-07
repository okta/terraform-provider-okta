resource "okta_group" "test_group" {
  name        = "testAcc_replace_with_uuid_group_updated"
  description = "Test group for updated resource change test"
}

resource "okta_resource_set" "test" {
  label       = "testAcc_replace_with_uuid"
  description = "A test resource set for updated resource change test"
  resources   = ["https://${var.hostname}/api/v1/groups/${okta_group.test_group.id}", "https://${var.hostname}/api/v1/users"]
}

data "okta_resource_set" "test" {
  id = okta_resource_set.test.id
}
