resource "okta_identity_source_bulk_group_memberships_upsert" "example" {
  identity_source_id = "<identity-source-id>"
  session_id = "<session-id>"

  # Optional fields
  # group_external_id = "<group-external-id>"
  # member_external_ids = "<member_external_ids>"
}
