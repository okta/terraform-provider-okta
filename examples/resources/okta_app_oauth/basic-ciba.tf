resource "okta_authenticator" "test1" {
  key                = "custom_app"
  name               = "VCRTestCustomAppAuth001"
  status             = "ACTIVE"
  agree_to_terms     = "true"
  legacy_ignore_name = false
  settings = jsonencode({
    "userVerification" : "REQUIRED",
    "appInstanceId" : "0oazl5jsgoLmguzkM1d7"
  })
  provider_json = jsonencode({
    "type" : "PUSH",
    "configuration" : {
      "fcm" : {
        "id" : "ppcrbgysxxBZDHPTv1d7"
      }
    }
  })
}

resource "okta_app_oauth" "test" {
  label = "GH2408_CIBA"
  type  = "web"
  grant_types = [
    "authorization_code",
    "refresh_token",
    "urn:openid:params:grant-type:ciba",
  ]
  redirect_uris  = ["https://localhost:4200/callback"]
  response_types = ["code"]

  backchannel_custom_authenticator_id = okta_authenticator.test1.id
}
