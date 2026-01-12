# Custom OTP with method-level settings
resource "okta_authenticator" "test" {
  name               = "Custom OTP Authenticator"
  key                = "custom_otp"
  status             = "ACTIVE"
  legacy_ignore_name = false
  settings = jsonencode({
    "protocol" : "TOTP",
    "encoding" : "base32",
    "algorithm" : "HMacSHA256",
    "timeIntervalInSeconds" : 30,
    "passCodeLength" : 6,
    "acceptableAdjacentIntervals" : 3
  })

  method {
    type   = "otp"
    status = "ACTIVE"
  }
}
