resource "okta_identity_source_session" "test" {
  identity_source_id = "0oaxc95befZNgrJl71d7"
}

resource "okta_identity_source_bulk_delete" "test" {
  identity_source_id = okta_identity_source_session.test.identity_source_id
  session_id         = okta_identity_source_session.test.id
  entity_type        = "USERS"

  profiles {
    external_id = "USEREXT123456TESTUSER1"
  }
}
