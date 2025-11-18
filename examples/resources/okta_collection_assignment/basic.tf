resource "okta_collection" "test" {
  name = "Test Collection"
}

resource "okta_group" "test" {
  name = "Test Group"
}

resource "okta_collection_assignment" "test" {
  collection_id   = okta_collection.test.id
  principal_id    = okta_group.test.id
  principal_type  = "OKTA_GROUP"
  actor          = "API"
}
