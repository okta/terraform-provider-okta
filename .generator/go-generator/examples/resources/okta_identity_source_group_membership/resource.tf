resource "okta_identity_source_group_membership" "example" {
  identity_source_id = "<identity-source-id>"
  group_or_external_id = "<group-or-external-id>"

  # Optional fields
  # member_external_ids = "<member_external_ids>"
  # member_external_id = "<member-external-id>"
}
