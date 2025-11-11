resource "okta_app_oauth" "test" {
  label                      = "testAcc_replace_with_uuid"
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["http://d.com/"]
  response_types             = ["code"]
  token_endpoint_auth_method = "client_secret_basic"
  consent_method             = "TRUSTED"
}

data "okta_app_oauth" "test" {
  id = okta_app_oauth.test.id
}

data "okta_app_oauth" "test_label" {
  label = okta_app_oauth.test.label
}
