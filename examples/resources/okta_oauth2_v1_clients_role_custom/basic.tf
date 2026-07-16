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
  role         = "cr0ohugfh3nITgFgm1d7"
  resource_set = "iamohu8okbUTSrARn1d7"
}
