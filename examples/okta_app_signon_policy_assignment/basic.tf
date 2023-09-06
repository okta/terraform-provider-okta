resource "okta_app_oauth" "test" {
  label                      = "testAcc_App_replace_with_uuid"
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["http://d.com/"]
  response_types             = ["code"]
  client_basic_secret        = "something_from_somewhere"
  client_id                  = "cid_replace_with_uuid"
  token_endpoint_auth_method = "client_secret_basic"
  consent_method             = "TRUSTED"
  wildcard_redirect          = "DISABLED"
}
resource "okta_app_signon_policy" "policy_1" {
  name        = "testAcc_SignOn_Policy_1_replace_with_uuid"
  description = "Policy 1"
}

resource "okta_app_signon_policy" "policy_2" {
  name        = "testAcc_SignOn_Policy_2_replace_with_uuid"
  description = "Policy 2"
}

resource "okta_app_signon_policy_assignment" "test" {
  app_id    = okta_app_oauth.test.id
  policy_id = okta_app_signon_policy.policy_1.id
}
