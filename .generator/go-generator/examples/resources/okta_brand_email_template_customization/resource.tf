resource "okta_brand_email_template_customization" "example" {
  brand_id = "<brand-id>"
  template_name = "Example Template Name"
  body = "<body>"
  language = "<language>"
  subject = "<subject>"

  # Optional fields
  # is_default = true
}
