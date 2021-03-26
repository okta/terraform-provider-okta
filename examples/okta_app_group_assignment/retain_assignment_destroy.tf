resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["implicit", "authorization_code"]
  redirect_uris  = ["http://d.com/"]
  response_types = ["code", "token", "id_token"]
  issuer_mode    = "ORG_URL"

  lifecycle {
    ignore_changes = ["users", "groups"]
  }
}

resource "okta_group" "test" {
  name = "testAcc_replace_with_uuid"
}
