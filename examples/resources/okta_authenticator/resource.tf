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

# Phone authenticator with method-level control
resource "okta_authenticator" "phone" {
  name = "Phone Authentication"
  key  = "phone_number"

  # Enable SMS method
  method {
    type   = "sms"
    status = "ACTIVE"
  }

  # Disable voice method
  method {
    type   = "voice"
    status = "INACTIVE"
  }
}

# Okta Verify with specific methods enabled
resource "okta_authenticator" "okta_verify" {
  name = "Okta Verify"
  key  = "okta_verify"

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
    status = "INACTIVE"
  }
}

# Custom OTP with method-level settings
resource "okta_authenticator" "custom_otp_with_method" {
  name               = "Custom OTP Authenticator"
  key                = "custom_otp"
  status             = "ACTIVE"
  legacy_ignore_name = false

  method {
    type   = "otp"
    status = "ACTIVE"
    settings = jsonencode({
      "protocol" : "TOTP",
      "encoding" : "base32",
      "algorithm" : "HMacSHA256",
      "timeIntervalInSeconds" : 30,
      "passCodeLength" : 6,
      "acceptableAdjacentIntervals" : 3
    })
  }
}
