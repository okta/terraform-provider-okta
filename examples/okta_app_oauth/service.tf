resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["implicit", "client_credentials"]
  redirect_uris  = ["http://test.com"]
}
