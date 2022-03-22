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


  // Since Okta forces us to create it with a redirect URI we have to ignore future changes, they will be detected as config drift.
  lifecycle {
    ignore_changes = [post_logout_redirect_uris]
  }
}

resource "okta_app_oauth_post_logout_redirect_uri" "test" {
  app_id = okta_app_oauth.test.id
  uri    = "https://www.example-updated.com"
}
