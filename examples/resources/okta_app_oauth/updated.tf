resource "okta_app_oauth" "test" {
  label                                = "testAcc_replace_with_uuid"
  status                               = "INACTIVE"
  type                                 = "browser"
  grant_types                          = ["implicit"]
  redirect_uris                        = ["http://d-*.com/aaa"]
  response_types                       = ["token", "id_token"]
  hide_ios                             = true
  hide_web                             = true
  auto_submit_toolbar                  = false
  wildcard_redirect                    = "SUBDOMAIN"
  participate_slo                      = true
  frontchannel_logout_uri              = "http://d-*.com/logout"
  frontchannel_logout_session_required = true
}
