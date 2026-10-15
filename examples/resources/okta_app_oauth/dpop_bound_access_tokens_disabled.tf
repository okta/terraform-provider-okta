resource "okta_app_oauth" "test" {
  label                      = "testAcc_replace_with_uuid"
  type                       = "web"
  grant_types                = ["authorization_code"]
  redirect_uris              = ["https://example.com/callback"]
  response_types             = ["code"]
  dpop_bound_access_tokens   = false
  skip_authentication_policy = true
}
