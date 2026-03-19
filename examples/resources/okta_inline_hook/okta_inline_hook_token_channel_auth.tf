resource "okta_inline_hook" "token_channel_auth" {
  name    = "Inline Hook Channel Auth"
  version = "1.0.0"
  type    = "com.okta.oauth2.tokens.transform"

  channel = {
    version = "1.0.0"
    uri     = "https://example.com/test"
    method  = "POST"
  }

  auth = {
    key   = "Authorization"
    type  = "HEADER"
    value = "secret"
  }
}

