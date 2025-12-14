# Okta Verify with specific methods enabled
resource "okta_authenticator" "okta_verify" {
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

  # Disable signed nonce with specific settings
  method {
    type   = "signed_nonce"
    status = "INACTIVE"
    settings = jsonencode({
      "algorithms" : ["ES256", "ES384", "RS256"],
      "keyProtection" : "HARDWARE"
    })
  }
}
