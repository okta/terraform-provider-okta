resource "okta_app_oauth" "test" {
  label               = "testAcc_replace_with_uuid"
  status              = "ACTIVE"
  type                = "browser"
  grant_types         = ["authorization_code", "refresh_token"]
  redirect_uris       = ["http://d.com/aaa"]
  response_types      = ["code"]
  hide_ios            = true
  hide_web            = true
  auto_submit_toolbar = false
}
