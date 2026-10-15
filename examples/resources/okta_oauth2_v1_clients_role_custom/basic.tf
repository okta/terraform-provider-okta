resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["client_credentials"]
  jwks_uri       = "https://example.com"
}

resource "okta_oauth2_v1_clients_role_custom" "test" {
  client_id    = okta_app_oauth.test.client_id
  type         = "CUSTOM"
  role         = "cr0120m8opgg3lywG1d8"
  resource_set = "iamqspdvrtb1i0JWi1d7"
}
