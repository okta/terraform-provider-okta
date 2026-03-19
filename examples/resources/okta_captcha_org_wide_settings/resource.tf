resource "okta_captcha" "example" {
  name       = "My CAPTCHA"
  type       = "HCAPTCHA"
  site_key   = "some_key"
  secret_key = "some_secret_key"
}

resource "okta_captcha_org_wide_settings" "example" {
  captcha_id  = okta_captcha.test.id
  enabled_for = ["SSR"]
}

### The following example disables org-wide CAPTCHA.

resource "okta_captcha" "example" {
  name       = "My CAPTCHA"
  type       = "HCAPTCHA"
  site_key   = "some_key"
  secret_key = "some_secret_key"
}

resource "okta_captcha_org_wide_settings" "example" {
}
