resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["client_credentials"]
  jwks_uri       = "https://example.com"
}

resource "okta_oauth2_v1_clients_role_api_access_management_admin" "test" {
  client_id = okta_app_oauth.test.client_id
  type      = "API_ACCESS_MANAGEMENT_ADMIN"
}
