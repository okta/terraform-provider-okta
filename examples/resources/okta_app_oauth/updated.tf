resource "okta_app_oauth" "test" {
  label                                = "testAcc_replace_with_uuid"
  status                               = "INACTIVE"
  type                                 = "browser"
  grant_types                          = ["implicit"]
  redirect_uris                        = ["https://*.example.com/callback"]
  response_types                       = ["token", "id_token"]
  hide_ios                             = true
  hide_web                             = true
  auto_submit_toolbar                  = false
  issuer_mode                          = "ORG_URL"
  wildcard_redirect                    = "SUBDOMAIN"
  participate_slo                      = true
  frontchannel_logout_uri              = "https://*.example.com/logout"
  frontchannel_logout_session_required = true
}
