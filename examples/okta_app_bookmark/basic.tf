resource "okta_group" "group" {
  name = "testAcc_replace_with_uuid"
}

resource "okta_app_bookmark" "test" {
  label                 = "testAcc_replace_with_uuid"
  url                   = "https://test.com"
  authentication_policy = "some-authentication-policy-id"
  groups                = [okta_group.group.id]
}
