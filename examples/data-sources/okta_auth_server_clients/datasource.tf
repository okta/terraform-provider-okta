# Create an OAuth application first
data "okta_auth_server_clients" "test" {
  token_id = ""
  auth_server_id = ""
  client_id = ""
}
