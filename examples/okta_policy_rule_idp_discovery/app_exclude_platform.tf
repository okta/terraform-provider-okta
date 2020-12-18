resource "okta_policy_rule_idp_discovery" "test" {
  policyid = data.okta_policy.test.id
  priority = 1
  name     = "testAcc_replace_with_uuid"
  idp_type = "OKTA"

  app_exclude {
    type = "APP"
    id   = okta_app_oauth.test.id
  }

  platform_include {
    type    = "MOBILE"
    os_type = "ANY"
  }
}

data "okta_policy" "test" {
  name = "Idp Discovery Policy"
  type = "IDP_DISCOVERY"
}

resource "okta_app_oauth" "test" {
  label                      = "testAcc_replace_with_uuid"
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["http://d.com/"]
  response_types             = ["code"]
  client_basic_secret        = "something_from_somewhere"
  custom_client_id           = "something_from_somewhere"
  token_endpoint_auth_method = "client_secret_basic"
}
