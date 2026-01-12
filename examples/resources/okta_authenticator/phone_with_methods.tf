# Phone authenticator with method-level control
resource "okta_authenticator" "test" {
  name = "Phone"
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
