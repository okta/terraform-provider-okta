resource "okta_identity_source_bulk_upsert" "example" {
  identity_source_id = "<identity-source-id>"
  session_id = "<session-id>"

  # Optional fields
  # entity_type = "<entity_type>"
  # external_id = "<external-id>"
  # email = "user@example.com"
  # first_name = "Example First Name"
  # home_address = "<home_address>"
}
