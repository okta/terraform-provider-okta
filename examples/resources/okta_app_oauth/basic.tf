resource "okta_app_oauth" "test" {
  label                      = "testAcc_replace_with_uuid"
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["https://example.com/callback"]
  response_types             = ["code"]
  token_endpoint_auth_method = "client_secret_basic"
  consent_method             = "TRUSTED"
  issuer_mode                = "ORG_URL"
  wildcard_redirect          = "DISABLED"
  participate_slo            = true
  frontchannel_logout_uri    = "https://example.com/logout"
  groups_claim {
    type  = "EXPRESSION"
    value = "aa"
    name  = "bb"
  }
}
