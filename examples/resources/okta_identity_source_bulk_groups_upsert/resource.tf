resource "okta_identity_source_session" "test" {
  identity_source_id = "0oaxc95befZNgrJl71d7"
}

resource "okta_identity_source_bulk_groups_upsert" "test" {
  identity_source_id = okta_identity_source_session.test.identity_source_id
  session_id         = okta_identity_source_session.test.id

  profiles {
    external_id = "GROUPEXT123456784C2IF"

    group_profile {
      display_name = "Test Group"
      description  = "Test group for identity source"
    }
  }
}
