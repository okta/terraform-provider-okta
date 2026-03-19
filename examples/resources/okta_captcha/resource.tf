resource "okta_captcha" "example" {
  name       = "My CAPTCHA"
  type       = "HCAPTCHA"
  site_key   = "some_key"
  secret_key = "some_secret_key"
}
