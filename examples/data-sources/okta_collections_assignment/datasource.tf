resource "okta_collection" "test" {
  name = "testAcc-ds-collection-assignment-replace_with_uuid"
}

resource "okta_collections_assignment" "test" {
  collection_id   = okta_collection.test.id
  expiration_time = "2027-12-31T23:59:59Z"
  time_zone       = "Asia/Kolkata"

  principal {
    external_id = "00uto13wvy2Y5RCv91d7"
    type        = "OKTA_USER"
  }
}

data "okta_collections_assignment" "test" {
  collection_id = okta_collection.test.id

  depends_on = [okta_collections_assignment.test]
}
