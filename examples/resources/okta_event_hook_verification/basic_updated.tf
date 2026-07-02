resource "okta_event_hook" "example" {
  name = "testAcc_replace_with_uuid"
  events = [
    "user.lifecycle.create",
    "user.lifecycle.delete.initiated",
  ]

  channel = {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://eo4afyqp3adkxpk.m.pipedream.net"
  }

  auth = {
    type  = "HEADER"
    key   = "Authorization"
    value = "value"
  }
}

resource "okta_event_hook_verification" "user_assigned" {
  event_hook_id = okta_event_hook.example.id
}
