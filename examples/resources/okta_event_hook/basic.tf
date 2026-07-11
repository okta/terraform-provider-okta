resource "okta_event_hook" "test" {
  name = "testAcc_replace_with_uuid"
  channel {
    type    = "HTTP"
    version = "1.0.0"
    config {
      uri = "https://example.com/test"
    }
  }
  events {
    type  = "EVENT_TYPE"
    items = ["user.lifecycle.create", "user.lifecycle.delete.initiated"]
  }
}
