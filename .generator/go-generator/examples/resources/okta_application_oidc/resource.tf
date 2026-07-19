resource "okta_application_oidc" "example" {
  sign_on_mode = "<sign_on_mode>"
  label = "Example Label"
  name = "Example Name"
  grant_types = "<grant_types>"
  mode = "<mode>"
  connection = "<connection>"
  rotation_type = "<rotation_type>"

  # Optional fields
  # error_redirect_url = "https://example.com"
  # login_redirect_url = "https://example.com"
  # self_service = true
  # auto_key_rotation = true
  # client_id = "<client-id>"
}
