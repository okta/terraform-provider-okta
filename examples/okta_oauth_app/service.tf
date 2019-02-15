resource "okta_oauth_app" "test" {
  label          = "testAcc_%[1]d"
  type           = "service"
  response_types = ["token"]
  grant_types    = ["implicit", "client_credentials"]
  redirect_uris  = ["http://test.com"]
}
