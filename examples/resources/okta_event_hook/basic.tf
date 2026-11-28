resource "okta_event_hook" "test" {
  name = "testAcc_replace_with_uuid"

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

    filter {
      type = "EXPRESSION_LANGUAGE"
      event_filter_map {
        event = "user.lifecycle.create"
        condition {
          expression = "event.target.?[type eq 'User'].size()>0"
        }
      }
    }
  }
}
