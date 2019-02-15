resource "okta_oauth_app" "test" {
  label         = "testAcc_%[1]d"
  type          = "native"
  grant_types   = ["authorization_code"]
  redirect_uris = ["http://d.com/"]
}
