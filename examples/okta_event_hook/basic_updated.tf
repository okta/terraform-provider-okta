resource "okta_event_hook" "test" {
  name   = "testAcc_replace_with_uuid"
  status = "INACTIVE"
  events = [
    "user.lifecycle.create",
    "user.lifecycle.delete.initiated",
    "user.account.update_profile",
  ]

  headers {
    key   = "x-test-header"
    value = "test stuff"
  }

  headers {
    key   = "x-another-header"
    value = "more test stuff"
  }

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/testUpdated"
  }

  auth = {
    type  = "HEADER"
    key   = "Authorization"
    value = "123"
  }
}
