// This would normally be in another repo if you were decentralizing redirect_uri settings
resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["authorization_code"]
  response_types = ["code"]

  // Okta requires at least one redirect URI to create an app
  redirect_uris = ["myapp://callback"]

  // Ignore redirect uris if you are going to manage them with the
  // okta_app_oauth_redirect_uri resource and not have change detection on the
  // app for that value.
  lifecycle {
    ignore_changes = [redirect_uris]
  }
}

resource "okta_app_oauth_redirect_uri" "test" {
  app_id = okta_app_oauth.test.id
  uri    = "http://google.com"
}
