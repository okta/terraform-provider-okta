resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["client_credentials"]
  jwks_uri       = "https://example.com"
}
