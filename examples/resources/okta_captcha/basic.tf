resource "okta_captcha" "test" {
  name       = "testAcc_replace_with_uuid"
  type       = "HCAPTCHA"
  site_key   = "random_key"
  secret_key = "random_key"
}
