resource "okta_app_oauth" "test_app" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["authorization_code"]
  response_types = ["code"]
  redirect_uris  = ["http://d.com/"]
}

resource "okta_app_oauth_api_scope" "test_app_scopes" {
  app_id = okta_app_oauth.test_app.id
  issuer = "https://your.okta.org"
  scopes = ["okta.users.read"]
}
