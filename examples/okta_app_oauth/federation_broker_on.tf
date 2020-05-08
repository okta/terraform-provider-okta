resource "okta_app_oauth" "test" {
  label                     = "testAcc_replace_with_uuid"
  type                      = "web"
  grant_types               = ["implicit", "authorization_code"]
  redirect_uris             = ["http://d.com/"]
  post_logout_redirect_uris = ["http://d.com/post"]
  login_uri                 = "http://test.com"
  response_types            = ["code", "token", "id_token"]
  consent_method            = "TRUSTED"
  implicit_assignment       = true
}
