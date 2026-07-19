resource "okta_application_auto_login" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  login_url = "https://example.com"

  # Optional fields
  # error_redirect_url = "https://example.com"
  # login_redirect_url = "https://example.com"
  # self_service = true
  # algorithm = "<algorithm>"
  # digest_algorithm = "<digest_algorithm>"
}
