resource "okta_event_hook" "test" {
  name = "testAcc_replace_with_uuid"
  events = [
    "group.user_membership.add",
    "user.lifecycle.create",
  ]

  filter {
    event     = "group.user_membership.add"
    condition = "event.target.?[type eq 'UserGroup'].size()>0 && event.target.?[displayName eq 'Marketing'].size()>0"
  }

  filter {
    event     = "user.lifecycle.create"
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
