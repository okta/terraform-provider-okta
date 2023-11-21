resource "okta_inline_hook" "test" {
  name    = "testAcc_replace_with_uuid"
  status  = "ACTIVE"
  type    = "com.okta.import.transform"
  version = "1.0.2"

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test1"
    method  = "POST"
  }
}
