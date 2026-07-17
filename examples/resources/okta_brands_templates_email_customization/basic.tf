resource "okta_brands_templates_email_customization" "test" {
  brand_id      = "replace_with_uuid"
  template_name = "replace_with_uuid"
  body          = "test-body"
  is_default    = false
  language      = "en"
  subject       = "test-subject"
}
