resource "okta_collection" "test" {
  name        = "testAcc-ds-collection-replace_with_uuid"
  description = "Collection read by a data source test"
}

data "okta_collection" "test" {
  id = okta_collection.test.id
}
