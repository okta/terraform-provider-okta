resource "okta_email_sender" "test" {
  from_name    = "testAcc_replace_with_uuid"
  from_address = "no-reply@example.com"
  subdomain    = "mail"
}
