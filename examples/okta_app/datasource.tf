resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code", "token", "id_token"]
  issuer_mode    = "ORG_URL"
  consent_method = "TRUSTED"
}

data "okta_app" "test" {
  label = okta_app_oauth.test.label
}

data "okta_app" "test2" {
  id = okta_app_oauth.test.id
}

data "okta_app" "test3" {
  label_prefix = okta_app_oauth.test.label
}
