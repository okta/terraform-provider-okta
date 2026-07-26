resource "okta_api_servers_authorization_server" "test" {
  api_server_id = "replace_with_uuid"
  issuer = "test-issuer"
  status = "INACTIVE"
}
