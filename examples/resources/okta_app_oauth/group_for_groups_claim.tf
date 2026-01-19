// A groups claim allows for groups to be injected into ID tokens or access tokens.
// A group whitelist allows for limiting the groups injected on a per-application basis to
// reduce the size of the tokens. In reality, you would include the group ID in the 
// whitelist but for testing purpose the name is used here.
// Okta's documentation has more information.
// https://developer.okta.com/docs/guides/create-token-with-groups-claim/create-groups-claim/

resource "okta_group" "whitelist_group" {
  name = "testAcc_replace_with_uuid"
}

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

  profile = <<JSON
   {
    "groups": {
      "whitelist": [
          "${okta_group.whitelist_group.name}"
        ]
    }
  }
JSON
}
