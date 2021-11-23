resource "okta_captcha" "test" {
  name       = "testAcc_replace_with_uuid"
  type       = "HCAPTCHA"
  site_key   = "random_key"
  secret_key = "random_key"
}

resource "okta_captcha_org_wide_settings" "test" {
  captcha_id  = okta_captcha.test.id
  enabled_for = ["SSR", "SSPR", "SIGN_IN"]
}
