resource "okta_application_browser_plugin" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  name = "Example Name"
  button_field = "<button_field>"
  password_field = "<password_field>"
  url = "https://example.com"
  username_field = "Example Username Field"

  # Optional fields
  # error_redirect_url = "https://example.com"
  # login_redirect_url = "https://example.com"
  # self_service = true
  # algorithm = "<algorithm>"
  # digest_algorithm = "<digest_algorithm>"
}
