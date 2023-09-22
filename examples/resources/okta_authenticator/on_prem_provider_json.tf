resource "okta_authenticator" "test" {
  name = "On-Prem MFA"
  key  = "onprem_mfa"
  provider_json = jsonencode(
    {
      "type" : "DEL_OATH",
      "configuration" : {
        "authPort" : 999,
        "userNameTemplate" : {
          "template" : "global.assign.userName.login"
        },
        "hostName" : "localhost",
        "sharedSecret" : "Sh4r3d s3cr3t"
      }
    }
  )
}