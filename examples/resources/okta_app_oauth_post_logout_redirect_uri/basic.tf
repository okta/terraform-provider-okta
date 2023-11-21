// This would normally be in another repo if you were decentralizing redirect_uri settings
resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["authorization_code"]
  response_types = ["code"]

  // Okta requires at least one redirect URI to create an app
  redirect_uris = ["myapp://callback"]

  // After logout, Okta redirects users to one of these URIs
  post_logout_redirect_uris = ["https://www.example.com"]

  // Ignore post logout redirect uris if you are going to manage them with the
  // okta_app_oauth_post_logout_redirect_uri resource and not have change
  // detection on the app for that value.
  lifecycle {
    ignore_changes = [post_logout_redirect_uris]
  }
}

resource "okta_app_oauth_post_logout_redirect_uri" "test" {
  app_id = okta_app_oauth.test.id
  uri    = "http://google.com"
}
