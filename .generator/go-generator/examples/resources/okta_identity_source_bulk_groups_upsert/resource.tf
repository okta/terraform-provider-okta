resource "okta_identity_source_bulk_groups_upsert" "example" {
  identity_source_id = "<identity-source-id>"
  session_id = "<session-id>"

  # Optional fields
  # external_id = "<external-id>"
  # description = "Example description"
  # display_name = "Example Display Name"
}
