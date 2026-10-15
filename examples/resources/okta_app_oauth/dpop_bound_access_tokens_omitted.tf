// dpop_bound_access_tokens is deliberately absent here. enduser_note forces an update API call so
// the test can prove the provider does not reset an already-enabled DPoP setting back to false.
resource "okta_app_oauth" "test" {
  label                      = "testAcc_replace_with_uuid"
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["https://example.com/callback"]
  response_types             = ["code"]
  enduser_note               = "dpop omitted update"
  skip_authentication_policy = true
}
