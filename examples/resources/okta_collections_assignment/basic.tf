resource "okta_collection" "test" {
  name        = "testAcc-collection-assignment-replace_with_uuid"
  description = "Collection for an assignment test"
}

resource "okta_collections_assignment" "test" {
  collection_id = okta_collection.test.id

  principal {
    external_id = "00uto13wvy2Y5RCv91d7"
    type        = "OKTA_USER"
  }
}
