resource "okta_app_signon_policy" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "Test policy for office365_client_include"
}

resource "okta_app_signon_policy_rule" "test" {
  policy_id                = okta_app_signon_policy.test.id
  name                     = "testAcc_replace_with_uuid"
  office365_client_include = ["WEB", "MODERN_AUTH"]
}
