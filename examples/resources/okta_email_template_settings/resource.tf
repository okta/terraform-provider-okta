resource "okta_email_template_settings" "example" {
  brand_id      = "<brand id>"
  template_name = "ForgotPassword"
  recipients    = "ADMINS_ONLY"
}