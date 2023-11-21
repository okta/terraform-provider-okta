resource "okta_captcha" "test" {
  name       = "testAcc_replace_with_uuid_updated"
  type       = "HCAPTCHA"
  site_key   = "random_key_updated"
  secret_key = "random_key"
}
