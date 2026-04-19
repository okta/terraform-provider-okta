resource "okta_identity_source_session" "test" {
  identity_source_id = "0oaxc95befZNgrJl71d7"
}

resource "okta_identity_source_bulk_group_memberships_upsert" "test" {
  identity_source_id = okta_identity_source_session.test.identity_source_id
  session_id         = okta_identity_source_session.test.id

  memberships {
    group_external_id   = "GROUPEXT123456784C2IF"
    member_external_ids = ["USEREXT123456TESTUSER3"]
  }
}
