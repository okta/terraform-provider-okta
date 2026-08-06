resource "okta_brands_templates_email_customization" "example" {
  brand_id      = "<brand_id>"
  template_name = "ForgotPassword"
  language      = "en"
  subject       = "Forgot Password"
  body          = "<html><body>Reset your password here: $${resetPasswordLink}</body></html>"
  is_default    = true
}
