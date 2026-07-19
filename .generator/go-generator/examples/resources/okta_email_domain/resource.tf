resource "okta_email_domain" "example" {
  display_name = "Example Display Name"
  user_name = "Example User Name"
  brand_id = "<brand-id>"

  # Optional fields
  # fqdn = "<fqdn>"
  # record_type = "<record_type>"
  # verification_value = "<verification_value>"
  # domain = "<domain>"
  # validation_status = "ACTIVE"
}
