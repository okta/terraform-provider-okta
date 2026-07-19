resource "okta_application_bookmark" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  name = "Example Name"
  url = "https://example.com"

  # Optional fields
  # error_redirect_url = "https://example.com"
  # login_redirect_url = "https://example.com"
  # self_service = true
  # kid = "<kid>"
  # rotation_mode = "<rotation_mode>"
}
