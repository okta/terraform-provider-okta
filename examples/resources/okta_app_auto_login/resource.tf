resource "okta_app_auto_login" "example" {
  label                = "Example App"
  sign_on_url          = "https://example.com/login.html"
  sign_on_redirect_url = "https://example.com"
  reveal_password      = true
  credentials_scheme   = "EDIT_USERNAME_AND_PASSWORD"
}

resource "okta_app_auto_login" "example" {
  label             = "Google Example App"
  status            = "ACTIVE"
  preconfigured_app = "google"
  app_settings_json = <<JSON
{
    "domain": "okta",
    "afwOnly": false
}
JSON
}
