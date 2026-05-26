resource "okta_identity_source_group_membership" "example" {
  identity_source_id   = "0oaxc95befZNgrJl71d7"
  group_or_external_id = "GRPEXT123456TESTGROUP1"
  member_external_id   = "USEREXT123456TESTUSER1"
}
