resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "blah"
  login      = "collections_assignment_test@example.com"
  email      = "collections_assignment_test@example.com"
}

resource "okta_collection" "test" {
  name        = "testAcc-collection-assignment-replace_with_uuid"
  description = "Collection for an assignment test"
}

resource "okta_collections_assignment" "test" {
  collection_id = okta_collection.test.id

  principal {
    external_id = okta_user.test.id
    type        = "OKTA_USER"
  }
}
