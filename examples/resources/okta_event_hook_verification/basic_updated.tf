resource "okta_event_hook" "example" {
  name = "testAcc_replace_with_uuid"

  channel {
    type    = "HTTP"
    version = "1.0.0"
    config {
      uri = "https://eo4afyqp3adkxpk.m.pipedream.net"
    }
  }

  events {
    type  = "EVENT_TYPE"
    items = ["user.lifecycle.create", "user.lifecycle.delete.initiated"]
  }
}

resource "okta_event_hook_verification" "user_assigned" {
  event_hook_id = okta_event_hook.example.id
}
