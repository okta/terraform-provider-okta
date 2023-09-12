data "okta_policy" "test" {
  name = "Any two factors"
  type = "ACCESS_POLICY"
}

resource "okta_app_signon_policy" "test" {
  name        = "testAcc_Policy_replace_with_uuid"
  description = "Sign On Policy"
  depends_on = [
    data.okta_policy.test
  ]
}

resource "okta_app_bookmark" "test" {
  label                 = "testAcc_replace_with_uuid"
  url                   = "https://test.com"
  authentication_policy = okta_app_signon_policy.test.id
}
