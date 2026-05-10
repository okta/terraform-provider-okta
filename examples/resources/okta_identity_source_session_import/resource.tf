resource "okta_identity_source_session" "test" {
  identity_source_id = "0oaxc95befZNgrJl71d7"
}

resource "okta_identity_source_bulk_upsert" "test" {
  identity_source_id = okta_identity_source_session.test.identity_source_id
  session_id         = okta_identity_source_session.test.id

  entity_type = "USERS"

  profiles {
    external_id = "USEREXT00000000001"

    profile {
      user_name  = "test.user@example.com"
      first_name = "Test"
      last_name  = "User"
      email      = "test.user@example.com"
    }
  }
}

resource "okta_identity_source_session_import" "test" {
  identity_source_id = okta_identity_source_session.test.identity_source_id
  session_id         = okta_identity_source_session.test.id

  depends_on = [okta_identity_source_bulk_upsert.test]
}
