resource "okta_authorization_server_claim" "example" {
  auth_server_id = "<auth-server-id>"

  # Optional fields
  # always_include_in_token = true
  # claim_type = "<claim_type>"
  # scopes = "<scopes>"
  # group_filter_type = "<group_filter_type>"
  # name = "Example Name"
}
