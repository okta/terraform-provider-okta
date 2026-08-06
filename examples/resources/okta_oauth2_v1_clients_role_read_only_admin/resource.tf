resource "okta_oauth2_v1_clients_role_read_only_admin" "example" {
  client_id = "<client_id>"
  type      = "READ_ONLY_ADMIN"
}
