resource "okta_auth_server_scope" "example" {
  auth_server_id   = "<auth server id>"
  metadata_publish = "NO_CLIENTS"
  name             = "example"
  consent          = "IMPLICIT"
}
