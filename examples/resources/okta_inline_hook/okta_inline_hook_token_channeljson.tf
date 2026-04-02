resource "okta_inline_hook" "token_channeljson" {
  name    = "Inline Hook Channel JSON"
  type    = "com.okta.oauth2.tokens.transform"
  version = "1.0.0"
  status  = "ACTIVE"
  channel_json = jsonencode({
    "version" : "1.0.0",
    "type" : "HTTP",
    "config" : {
      "method" : "POST",
      "uri" : "https://httpbin.org/post",
      "headers" : [
        {
          "key" : "custom_key",
          "value" : "custom_val"
        }
      ],
      "authScheme" : {
        "type" : "HEADER",
        "key" : "authorization",
        "value" : "some_secret"
      }
    }
  })
}

