resource "okta_event_hook" "example" {
  name = "example"

  channel {
    type    = "HTTP"
    version = "1.0.0"
    config {
      uri = "https://example.com/test"
      auth_scheme {
        type  = "HEADER"
        key   = "Authorization"
        value = "123"
      }
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

resource "okta_event_hook_verification" "example" {
  event_hook_id = okta_event_hook.example.id
}
