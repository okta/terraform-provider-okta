resource "okta_event_hook" "test" {
  name = "testAcc_replace_with_uuid"
  events = [
    "user.lifecycle.create",
    "user.lifecycle.delete.initiated",
  ]

  filter {
    event     = "user.lifecycle.create"
    condition = "event.actor.id != null"
  }

  filter {
    event     = "user.lifecycle.delete.initiated"
    condition = "event.actor.id != null"
  }

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test-updated"
  }

  auth = {
    type  = "HEADER"
    key   = "Authorization"
    value = "123"
  }
}
