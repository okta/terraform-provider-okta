resource "okta_authenticator" "test" {
  name = "Security Question"
  key  = "security_question"
  settings = jsonencode(
    {
      "allowedFor" : "recovery"
    }
  )
}

resource "okta_authenticator" "otp" {
  name   = "Custom OTP"
  key    = "custom_otp"
  status = "ACTIVE"
  settings = jsonencode({
    "protocol" : "TOTP",
    "acceptableAdjacentIntervals" : 3,
    "timeIntervalInSeconds" : 30,
    "encoding" : "base32",
    "algorithm" : "HMacSHA256",
    "passCodeLength" : 6
  })
  // required to be false for custom_otp
  legacy_ignore_name = false
}
