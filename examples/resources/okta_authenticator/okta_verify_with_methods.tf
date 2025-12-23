# Okta Verify with specific methods enabled
resource "okta_authenticator" "test" {
  name = "Okta Verify"
  key  = "okta_verify"

  # Enable push notifications
  method {
    type   = "push"
    status = "ACTIVE"
  }

  # Enable TOTP
  method {
    type   = "totp"
    status = "ACTIVE"
  }

  # Enable signed nonce with specific settings
  method {
    type   = "signed_nonce"
    status = "ACTIVE"
    settings = jsonencode({
      "algorithms" : ["ES256", "RS256"],
      "keyProtection" : "ANY"
    })
  }
}
