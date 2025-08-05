resource "okta_inline_hook" "test" {
  name         = "testAcc_replace_with_uuid_%s"
  type         = "com.okta.saml.tokens.transform"
  version      = "1.0.2"
  status       = "ACTIVE"
  channel_json = <<JSON
{
        "type": "OAUTH",
        "version": "1.0.0",
        "config": {
            "headers": [
                {
                    "key": "Field 1",
                    "value": "Value 1"
                },
                {
                    "key": "Field 2",
                    "value": "Value 2"
                }
            ],
            "method": "POST",
            "authType": "client_secret_post",
            "uri": "https://example.com/service",
            "clientId": "abc123",
            "clientSecret": "def456_UPDATED",
            "tokenUrl": "https://example.com/token",
            "scope": "api"
        }
}
JSON
}
 