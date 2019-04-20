resource "okta_oauth_app" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code", "token", "id_token"]
  issuer_mode    = "ORG_URL"
}

data "okta_app" "test" {
  label = "${okta_oauth_app.test.label}"
}

data "okta_app" "test2" {
  id = "${okta_oauth_app.test.id}"
}

data "okta_app" "test3" {
  label_prefix = "${okta_oauth_app.test.label}"
}
