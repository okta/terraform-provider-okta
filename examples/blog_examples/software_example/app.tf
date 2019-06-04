data okta_group everyone {
  name = "Everyone"
}

resource okta_oauth_app my_app {
  label = "My App"
  type  = "web"

  grant_types = [
    "authorization_code",
    "refresh_token",
    "implicit",
  ]

  response_types = [
    "id_token",
    "code",
  ]

  redirect_uris             = ["https://example.com/auth-callback"]
  post_logout_redirect_uris = ["https://example.com"]
  login_uri                 = "https://example.com"
  groups                    = ["${data.okta_group.everyone.id}"]
}
