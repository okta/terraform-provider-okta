resource "okta_event_hook" "example" {
  name = "example"
  events = [
    "user.lifecycle.create",
    "user.lifecycle.delete.initiated",
  ]

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/test"
  }

  auth = {
    type  = "HEADER"
    key   = "Authorization"
    value = "123"
  }

  filter {
    event = "user.lifecycle.create"
    condition {
      expression = "event.actor.id eq '1234'"
    }
  }
}