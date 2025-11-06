resource "okta_authenticator" "okta_verify" {
  legacy_ignore_name = false
  key                = "okta_verify"
  name               = "Okta Verify"
  settings = jsonencode({
    "channelBinding" : {
      "required" : "ALWAYS",
    "style" : "NUMBER_CHALLENGE" },
    "compliance" : { "fips" : "OPTIONAL" },
    "userVerification" : "PREFERRED",
    "enrollmentSecurityLevel" : "HIGH",
    "userVerificationMethods" : ["BIOMETRICS"]
  })
  status = "ACTIVE"
}

