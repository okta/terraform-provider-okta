resource "okta_captcha" "test" {
  name       = "testAcc_replace_with_uuid"
  type       = "HCAPTCHA"
  site_key   = "random_key"
  secret_key = "random_secret_key"
}

resource "okta_org_captcha" "test" {
  captcha_id    = okta_captcha.test.id
  enabled_pages = ["SIGN_IN"]
}

data "okta_org_captcha" "test" {
  depends_on = [okta_org_captcha.test]
}
