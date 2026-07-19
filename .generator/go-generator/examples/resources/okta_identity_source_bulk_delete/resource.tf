resource "okta_identity_source_bulk_delete" "example" {
  identity_source_id = "<identity-source-id>"
  session_id = "<session-id>"

  # Optional fields
  # entity_type = "<entity_type>"
  # external_id = "<external-id>"
}
