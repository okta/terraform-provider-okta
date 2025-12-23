resource "okta_authenticator" "test" {
  name   = "Okta Verify"
  key    = "okta_verify"
  status = "ACTIVE"
  settings = jsonencode({
    "userVerification" : "PREFERRED",
    "channelBinding" : {
      "required" : "ALWAYS",
      "style" : "NUMBER_CHALLENGE"
    }
  })

  method {
    type   = "push"
    status = "ACTIVE"
  }

  method {
    type   = "totp"
    status = "ACTIVE"
  }

  method {
    type   = "signed_nonce"
    status = "ACTIVE"
  }
}
