resource "okta_event_hook" "example" {
  name        = "example"
  description = "Notifies an external service of user lifecycle events"

  channel {
    type    = "HTTP"
    version = "1.0.0"
    config {
      uri = "https://example.com/test"

      auth_scheme {
        type  = "HEADER"
        key   = "Authorization"
        value = "my-shared-secret"
      }

      headers {
        key   = "X-Other-Header"
        value = "my-header-value"
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
