resource "okta_application_saml_11" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  name = "Example Name"

  # Optional fields
  # error_redirect_url = "https://example.com"
  # login_redirect_url = "https://example.com"
  # self_service = true
  # kid = "<kid>"
  # rotation_mode = "<rotation_mode>"
}
