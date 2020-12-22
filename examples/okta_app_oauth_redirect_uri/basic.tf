// This would normally be in another repo if you were decentralizing redirect_uri settings
resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["authorization_code"]
  response_types = ["code"]

  // Okta requires at least one redirect URI to create an app
  redirect_uris = ["myapp://callback"]

  // Since Okta forces us to create it with a redirect URI we have to ignore future changes, they will be detected as config drift.
  lifecycle {
    ignore_changes = ["redirect_uris"]
  }
}

resource "okta_app_oauth_redirect_uri" "test" {
  app_id = okta_app_oauth.test.id
  uri    = "http://google.com"
}
