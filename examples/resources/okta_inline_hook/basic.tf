resource "okta_inline_hook" "test" {
  name    = "testAcc_replace_with_uuid"
  type    = "com.okta.oauth2.tokens.transform"
  version = "1.0.0"
  channel {
    type    = "HTTP"
    version = "1.0.0"
  }
}
