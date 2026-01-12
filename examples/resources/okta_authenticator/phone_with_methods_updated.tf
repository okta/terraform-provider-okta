resource "okta_authenticator" "test" {
  name   = "Phone"
  key    = "phone_number"
  status = "ACTIVE"

  method {
    type   = "sms"
    status = "ACTIVE"
  }

  method {
    type   = "voice"
    status = "ACTIVE"
  }
}
