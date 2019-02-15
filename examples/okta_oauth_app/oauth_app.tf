resource "okta_oauth_app" "test" {
  label          = "testAcc_%[1]d"
  type           = "web"
  grant_types    = ["authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code"]
}
