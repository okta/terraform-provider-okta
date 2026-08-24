resource "okta_captcha" "test" {
  name       = "testAcc_replace_with_uuid"
  type       = "HCAPTCHA"
  site_key   = "random_key"
  secret_key = "random_secret_key"
}

data "okta_captcha" "test" {
  id = okta_captcha.test.id
}
