resource "okta_email_sender" "example" {
  from_name    = "Paul Atreides"
  from_address = "no-reply@caladan.planet"
  subdomain    = "mail"
}

resource "okta_email_sender_verification" "example" {
  sender_id = okta_email_sender.valid.id
}
