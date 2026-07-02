resource "okta_identity_source_import" "test" {
  identity_source_id = "0oaxc95befZNgrJl71d7"

  upsert_users {
    entity_type = "USERS"

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
}
