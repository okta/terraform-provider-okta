resource "okta_email_sender" "example" {
  from_name    = "Paul Atreides"
  from_address = "no-reply@caladan.planet"
  subdomain    = "mail"
}
