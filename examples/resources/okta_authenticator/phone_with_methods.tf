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
