resource "okta_authenticator" "test1" {
  key                = "custom_app"
  name               = "VCRTestCustomAppAuthNewRenamed"
  status             = "ACTIVE"
  agree_to_terms     = "true"
  legacy_ignore_name = false
  settings = jsonencode({
    "userVerification" : "PREFERRED",
    "appInstanceId" : "0oa12345678ABCDEFG"
  })
  provider_json = jsonencode({
    "type" : "PUSH",
    "configuration" : {
      "fcm" : {
        "id" : "ppcABCDEFGH12345678"
      }
    }
  })
}
