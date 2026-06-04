resource "okta_authenticator" "test1" {
  key                = "custom_app"
  name               = "testAcc_replace_with_uuid"
  status             = "ACTIVE"
  agree_to_terms     = "true"
  legacy_ignore_name = false
  settings = jsonencode({
    "userVerification" : "REQUIRED",
    "appInstanceId" : "0oabcdefghi123456789"
  })
  provider_json = jsonencode({
    "type" : "PUSH",
    "configuration" : {
      "fcm" : {
        "id" : "ppcabcdefghijklmno123"
      }
    }
  })
}

resource "okta_app_oauth" "test" {
  label = "testAcc_replace_with_uuid"
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
