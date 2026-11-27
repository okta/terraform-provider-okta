resource "okta_event_hook" "test" {
  name        = "testAcc_replace_with_uuid"
  description = "testAcc_replace_with_uuid description"

  channel {
    type    = "HTTP"
    version = "1.0.0"
    config {
      uri = "https://example.com/testAcc_replace_with_uuid"
    }
  }

  events {
    type = "EVENT_TYPE"
    items = [
      "user.lifecycle.create",
      "user.lifecycle.delete.initiated",
    ]
  }
}

data "okta_event_hook" "test" {
  id = okta_event_hook.test.id
}
