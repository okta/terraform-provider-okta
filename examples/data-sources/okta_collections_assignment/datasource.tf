resource "okta_collection" "test" {
  name = "testAcc-ds-collection-assignment-replace_with_uuid"
}

resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "blah"
  login      = "collections_assignment_test@example.com"
  email      = "collections_assignment_test@example.com"
}


resource "okta_collections_assignment" "test" {
  collection_id   = okta_collection.test.id
  expiration_time = "2027-12-31T23:59:59Z"
  time_zone       = "Asia/Kolkata"

  principal {
    external_id = okta_user.test.id
    type        = "OKTA_USER"
  }
}

data "okta_collections_assignment" "test" {
  collection_id = okta_collection.test.id

  depends_on = [okta_collections_assignment.test]
}
