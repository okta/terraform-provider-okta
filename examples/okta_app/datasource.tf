resource "okta_oauth_app" "testAcc_%[1]d" {
  label          = "testAcc_%[1]d"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code", "token", "id_token"]
}

data "okta_app" "test" {
  label = "testAcc_%[1]d"
}

data "okta_app" "test2" {
  id = "${okta_oauth_app.testAcc_%[1]d.id}"
}
