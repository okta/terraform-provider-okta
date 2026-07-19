resource "okta_application_ws_federation" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  name = "Example Name"
  audience_restriction = "<audience_restriction>"
  authn_context_class_ref = "<authn_context_class_ref>"
  group_value_format = "<group_value_format>"
  name_id_format = "Example Name Id Format"
  site_url = "https://example.com"
  username_attribute = "Example Username Attribute"
  w_reply_url = "https://example.com"

  # Optional fields
  # error_redirect_url = "https://example.com"
  # login_redirect_url = "https://example.com"
  # self_service = true
  # kid = "<kid>"
  # rotation_mode = "<rotation_mode>"
}
