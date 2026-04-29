resource "okta_identity_source_session" "test" {
  identity_source_id = "0oaxc95befZNgrJl71d7"
}

resource "okta_identity_source_bulk_upsert" "test" {
  identity_source_id = okta_identity_source_session.test.identity_source_id
  session_id         = okta_identity_source_session.test.id
  entity_type        = "USERS"

  profiles {
    external_id = "USEREXT123456TESTUSER3"

    profile {
      user_name  = "testuser3@example.com"
      email      = "testuser3@example.com"
      first_name = "Test"
      last_name  = "User3"
    }
  }
}
