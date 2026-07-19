resource "okta_identity_source_bulk_groups_delete" "example" {
  identity_source_id = "<identity-source-id>"
  session_id = "<session-id>"

  # Optional fields
  # external_ids = "<external_ids>"
}
